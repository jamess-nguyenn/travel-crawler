package accommodation

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
	Get(
		ctx context.Context,
		req GetRequest,
	) (*Accommodation, error)

	List(
		ctx context.Context,
		req ListRequest,
	) (*response.ListResponse[Accommodation], error)
}

type SyncService interface {
	Sync(
		ctx context.Context,
		providerCode string,
		req SearchRequest,
	) error
}

type service struct {
	txManager         *database.TxManager
	accommodationRepo AccommodationRepository
	logger            *zap.SugaredLogger
}

type syncService struct {
	txManager         *database.TxManager
	accommodationRepo AccommodationRepository
	providers         ProviderManager
	logger            *zap.SugaredLogger
}

func NewService(
	txManager *database.TxManager,
	accommodationRepo AccommodationRepository,
	logger *zap.SugaredLogger,
) Service {
	return &service{
		txManager:         txManager,
		accommodationRepo: accommodationRepo,
		logger:            logger,
	}
}

func NewSyncService(
	txManager *database.TxManager,
	accommodationRepo AccommodationRepository,
	providers ProviderManager,
	logger *zap.SugaredLogger,
) SyncService {
	return &syncService{
		txManager:         txManager,
		accommodationRepo: accommodationRepo,
		providers:         providers,
		logger:            logger,
	}
}

func (s *service) Get(
	ctx context.Context,
	req GetRequest,
) (*Accommodation, error) {
	logger := applog.FromContext(s.logger, ctx)

	conditions := map[string]any{
		"id": req.ID,
	}

	accommodation, err := s.accommodationRepo.Get(ctx, conditions)

	getLogger := applog.WithGet(
		logger,
		conditions,
		accommodation,
	)

	if err != nil {
		getLogger.Errorw(
			"Fail to get accommodation",
			"error", err,
		)

		return nil, err
	}

	getLogger.Info("Get accommodation successfully")

	return accommodation, nil
}

func (s *service) buildListQuery(
	req ListRequest,
) database.ListQuery {
	conditions := make(map[string]any)

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
) (*response.ListResponse[Accommodation], error) {
	logger := applog.FromContext(s.logger, ctx)

	listQuery := s.buildListQuery(req)

	accommodations, total, err := s.accommodationRepo.List(ctx, listQuery)

	listLogger := applog.WithList(
		logger,
		listQuery,
		total,
		accommodations,
	)

	if err != nil {
		listLogger.Errorw(
			"Fail to list accommodations",
			"error", err,
		)

		return nil, err
	}

	listLogger.Info("List accommodations successfully")

	return &response.ListResponse[Accommodation]{
		Items:      accommodations,
		Pagination: response.NewPagination(req.GetPage(), req.GetPageSize(), total),
		Request:    req,
	}, nil
}

func (s *syncService) Sync(
	ctx context.Context,
	providerCode string,
	req SearchRequest,
) error {
	logger := applog.FromContext(s.logger, ctx)

	p, err := s.providers.Get(providerCode)
	if err != nil {
		logger.Warnw(
			"Fail to get provider",
			"provider_code", providerCode,
			"error", err,
		)

		return err
	}

	req.PageSize = 1000

	for page := 1; ; page++ {
		// Keep the original request intact except for the page.
		pageReq := req
		pageReq.Page = page

		resp, err := p.Search(ctx, pageReq)
		if err != nil {
			logger.Errorw(
				"Fail to search accommodations",
				"provider_code", providerCode,
				"request", pageReq,
				"error", err,
			)

			return fmt.Errorf(
				"search accommodations from %s, page %d: %w",
				providerCode,
				page,
				err,
			)
		}

		accommodations := resp.Accommodations
		total := len(accommodations)
		offset := helper.CalculateOffset(page, pageReq.PageSize)

		pageLogger := logger.With(
			"provider_code", providerCode,
			"request", pageReq,
			"offset", offset,
			"total", total,
			"data", accommodations,
		)

		pageLogger.Info("Search accommodations successfully")

		if total == 0 {
			pageLogger.Info("No more accommodations to sync")
			break
		}

		if err := s.persistAccommodation(ctx, accommodations, logger); err != nil {
			return fmt.Errorf(
				"persist accommodations from %s, page %d: %w",
				providerCode,
				page,
				err,
			)
		}

		pageLogger.Info("Sync accommodation page successfully")

		// Prefer the provider's explicit HasMore value.
		if hasMore, known := resp.Pagination.HasMoreValue(); known {
			if !hasMore {
				break
			}

			continue
		}

		// Otherwise, use the total number of matching records.
		if totalRecord, known := resp.Pagination.TotalValue(); known {
			if offset+total >= totalRecord {
				break
			}

			continue
		}

		// If pagination metadata is unavailable, assume that a
		// full page may have another page.
		if total < pageReq.PageSize {
			break
		}
	}

	return nil
}

func (s *syncService) persistAccommodation(
	ctx context.Context,
	accommodations []Accommodation,
	logger *zap.SugaredLogger,
) error {
	options := database.UpsertOptions{
		ConflictColumns: UpsertConflictColumns,
		UpdateColumns:   UpsertUpdateColumns,
	}

	batches := helper.Chunk(accommodations, constant.DatabaseBatchSize)
	totalBatch := len(batches)

	for i, batch := range batches {
		batchNumber := i + 1
		batchProgress := fmt.Sprintf("%d/%d", batchNumber, totalBatch)

		batchLogger := applog.WithBatch(
			logger,
			len(accommodations),
			batchProgress,
			len(batch),
			batch,
		)

		if err := s.accommodationRepo.Upsert(ctx, &batch, options); err != nil {
			batchLogger.Errorw(
				"Fail to upsert accommodation batch",
				"error", err,
			)

			return fmt.Errorf("upsert accommodation batch %s: %w", batchProgress, err)
		}

		batchLogger.Info("Upsert accommodation batch successfully")
	}

	return nil
}
