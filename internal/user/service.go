package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/database/scope"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/response"

	"travel-crawler/internal/database"
)

type Service interface {
	Upsert(
		ctx context.Context,
		status Status,
	) error

	Get(
		ctx context.Context,
		req GetRequest,
	) (*User, error)

	UpdateByID(
		ctx context.Context,
		ID uint64,
		actorID uint64,
		req UpdateRequest,
	) error

	List(
		ctx context.Context,
		req ListRequest,
	) (*response.ListResponse[User], error)

	ListAll(
		ctx context.Context,
		req ListAllRequest,
	) ([]User, error)
}

type service struct {
	txManager *database.TxManager
	userRepo  UserRepository
	logger    *zap.SugaredLogger
}

func NewService(
	txManager *database.TxManager,
	userRepo UserRepository,
	logger *zap.SugaredLogger,
) Service {
	return &service{
		txManager: txManager,
		userRepo:  userRepo,
		logger:    logger,
	}
}

func (s *service) Upsert(
	ctx context.Context,
	status Status,
) error {
	var users []User

	if err := s.persistUser(ctx, users); err != nil {
		return err
	}

	return nil
}

func (s *service) persistUser(
	ctx context.Context,
	users []User,
) error {
	logger := applog.FromContext(s.logger, ctx)

	options := database.UpsertOptions{
		ConflictColumns: UpsertConflictColumns,
		UpdateColumns:   UpsertUpdateColumns,
	}

	batches := helper.Chunk(users, constant.DatabaseBatchSize)
	totalBatch := len(batches)

	for i, batch := range batches {
		batchNumber := i + 1
		batchProgress := fmt.Sprintf("%d/%d", batchNumber, totalBatch)

		batchLogger := applog.WithBatch(
			logger,
			len(users),
			batchProgress,
			len(batch),
			batch,
		)

		if err := s.userRepo.Upsert(ctx, &batch, options); err != nil {
			batchLogger.Errorw(
				"Fail to upsert user batch",
				"error", err,
			)

			return fmt.Errorf("upsert user batch %s: %w", batchProgress, err)
		}

		batchLogger.Info("Upsert user batch successfully")
	}

	return nil
}

func (s *service) get(
	ctx context.Context,
	conditions map[string]any,
) (*User, error) {
	logger := applog.FromContext(s.logger, ctx)

	user, err := s.userRepo.Get(ctx, conditions)

	getLogger := applog.WithGet(
		logger,
		conditions,
		user,
	)

	if err != nil {
		getLogger.Errorw(
			"Fail to get user",
			"error", err,
		)

		return nil, err
	}

	getLogger.Info("Get user successfully")

	return user, nil
}

func (s *service) Get(
	ctx context.Context,
	req GetRequest,
) (*User, error) {
	conditions := map[string]any{
		"id": req.ID,
	}

	return s.get(ctx, conditions)
}

func (s *service) UpdateByID(
	ctx context.Context,
	ID uint64,
	actorID uint64,
	req UpdateRequest,
) error {
	logger := applog.FromContext(s.logger, ctx)

	conditions := map[string]any{
		"id": ID,
	}

	if _, err := s.get(ctx, conditions); err != nil {
		return err
	}

	now := helper.NowUTC()

	updates := map[string]any{
		"first_name":   req.FirstName,
		"last_name":    req.LastName,
		"full_name":    BuildFullName(req.FirstName, req.LastName),
		"phone_number": req.PhoneNumber,
		"status":       req.Status,
		"updated_at":   now,
		"updated_by":   actorID,
	}

	updateLogger := applog.WithUpdate(
		logger,
		ID,
		updates,
	)

	if err := s.userRepo.UpdateByID(ctx, ID, updates); err != nil {
		updateLogger.Errorw(
			"Fail to update user",
			"error", err,
		)

		return err
	}

	updateLogger.Info("Update user successfully")

	return nil
}

func (s *service) buildListQuery(
	req ListRequest,
) database.ListQuery {
	conditions := make(map[string]any)

	if req.Status != nil && *req.Status != "" {
		conditions["status"] = *req.Status
	}

	filters := []scope.Scope{
		scope.SearchColumns([]string{"email", "full_name", "phone_number"}, req.Search),
		scope.TimeRange("created_at", req.DateFrom, req.DateTo),
	}

	return database.ListQuery{
		Conditions: conditions,
		Filters:    filters,
		Scopes:     scope.BuildListScopes(req, sortColumns, defaultSortColumn),
		Columns:    nil,
	}
}

func (s *service) List(
	ctx context.Context,
	req ListRequest,
) (*response.ListResponse[User], error) {
	logger := applog.FromContext(s.logger, ctx)

	listQuery := s.buildListQuery(req)

	users, total, err := s.userRepo.List(ctx, listQuery)

	listLogger := applog.WithList(
		logger,
		listQuery,
		total,
		users,
	)

	if err != nil {
		listLogger.Errorw(
			"Fail to list users",
			"error", err,
		)

		return nil, err
	}

	listLogger.Info("List users successfully")

	return &response.ListResponse[User]{
		Items:      users,
		Pagination: response.NewPagination(req.GetPage(), req.GetPageSize(), total),
		Request:    req,
	}, nil
}

func (s *service) buildListAllQuery(
	req ListAllRequest,
) database.ListQuery {
	conditions := make(map[string]any)

	if req.Status != nil && *req.Status != "" {
		conditions["status"] = *req.Status
	}

	filters := make([]scope.Scope, 0)

	if len(req.IDs) > 0 {
		filters = append(filters, scope.WhereIn("id", req.IDs))
	}

	if len(req.Emails) > 0 {
		filters = append(filters, scope.WhereIn("email", req.Emails))
	}

	return database.ListQuery{
		Conditions: conditions,
		Filters:    filters,
		Scopes:     []scope.Scope{scope.OrderByDesc(defaultSortColumn)},
		Columns:    nil,
	}
}

func (s *service) ListAll(
	ctx context.Context,
	req ListAllRequest,
) ([]User, error) {
	logger := applog.FromContext(s.logger, ctx)

	listQuery := s.buildListAllQuery(req)

	users, err := s.userRepo.ListAll(ctx, listQuery)

	listLogger := applog.WithList(
		logger,
		listQuery,
		int64(len(users)),
		users,
	)

	if err != nil {
		listLogger.Errorw(
			"Fail to list all users",
			"error", err,
		)

		return nil, err
	}

	listLogger.Info("List all users successfully")

	return users, nil
}
