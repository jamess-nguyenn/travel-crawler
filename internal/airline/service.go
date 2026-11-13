package airline

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
	) (*Airline, error)

	List(
		ctx context.Context,
		req ListRequest,
	) (*response.ListResponse[Airline], error)
}

type SyncService interface {
	Sync(
		ctx context.Context,
		providerCode string,
		req SearchRequest,
	) error
}

type service struct {
	txManager   *database.TxManager
	airlineRepo AirlineRepository
	logger      *zap.SugaredLogger
}

type syncService struct {
	txManager   *database.TxManager
	airlineRepo AirlineRepository
	providers   ProviderManager
	logger      *zap.SugaredLogger
}

func NewService(
	txManager *database.TxManager,
	airlineRepo AirlineRepository,
	logger *zap.SugaredLogger,
) Service {
	return &service{
		txManager:   txManager,
		airlineRepo: airlineRepo,
		logger:      logger,
	}
}

func NewSyncService(
	txManager *database.TxManager,
	airlineRepo AirlineRepository,
	providers ProviderManager,
	logger *zap.SugaredLogger,
) SyncService {
	return &syncService{
		txManager:   txManager,
		airlineRepo: airlineRepo,
		providers:   providers,
		logger:      logger,
	}
}

func (s *service) Get(
	ctx context.Context,
	req GetRequest,
) (*Airline, error) {
	logger := applog.FromContext(s.logger, ctx)

	conditions := map[string]any{"id": req.ID}

	airline, err := s.airlineRepo.Get(ctx, conditions, selectColumns)

	getLogger := applog.WithGet(
		logger,
		conditions,
		selectColumns,
		airline,
	)

	if err != nil {
		getLogger.Errorw(
			"Fail to get airline",
			"error", err,
		)

		return nil, err
	}

	getLogger.Info("Get airline successfully")

	return airline, nil
}

func (s *service) buildListQuery(
	req ListRequest,
) database.ListQuery {
	conditions := make(map[string]any)

	if helper.IsNonEmptyString(req.Alliance) {
		conditions["alliance"] = *req.Alliance
	}

	filters := []scope.Scope{
		scope.SearchColumns([]string{"iata_code", "name"}, req.Search),
		scope.TimeRange("created_at", req.DateFrom, req.DateTo),
	}

	return database.ListQuery{
		Conditions: conditions,
		Filters:    filters,
		Scopes:     scope.BuildListScopes(req, sortColumns, defaultSortColumn),
		Columns:    selectColumns,
	}
}

func (s *service) List(
	ctx context.Context,
	req ListRequest,
) (*response.ListResponse[Airline], error) {
	logger := applog.FromContext(s.logger, ctx)

	listQuery := s.buildListQuery(req)

	airlines, total, err := s.airlineRepo.List(ctx, listQuery)

	listLogger := applog.WithList(
		logger,
		listQuery,
		total,
		airlines,
	)

	if err != nil {
		listLogger.Errorw(
			"Fail to list airlines",
			"error", err,
		)

		return nil, err
	}

	listLogger.Info("List airlines successfully")

	return &response.ListResponse[Airline]{
		Items:      airlines,
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

	req.PageSize = 10000

	resp, err := p.Search(ctx, req)
	if err != nil {
		logger.Errorw(
			"Fail to search airlines",
			"provider_code", providerCode,
			"request", req,
			"error", err,
		)

		return fmt.Errorf(
			"search airlines from %s: %w",
			providerCode,
			err,
		)
	}

	airlines := resp.Airlines
	total := len(airlines)

	listLogger := logger.With(
		"provider_code", providerCode,
		"request", req,
		"total", total,
		"data", helper.Sample(airlines, constant.SampleLogSize),
	)

	listLogger.Info("Search airlines successfully")

	if total == 0 {
		listLogger.Info("No airlines to sync")
		return nil
	}

	if err := s.persistAirline(ctx, airlines, logger); err != nil {
		return fmt.Errorf(
			"persist airlines from %s: %w",
			providerCode,
			err,
		)
	}

	listLogger.Info("Sync airlines successfully")

	return nil
}

func (s *syncService) persistAirline(
	ctx context.Context,
	airlines []Airline,
	logger *zap.SugaredLogger,
) error {
	total := len(airlines)
	if total == 0 {
		return nil
	}

	options := database.UpsertOptions{
		ConflictColumns: upsertConflictColumns,
		UpdateColumns:   upsertUpdateColumns,
	}

	batches := helper.Chunk(airlines, constant.DatabaseBatchSize)
	totalBatch := len(batches)

	for i, batch := range batches {
		batchNumber := i + 1
		batchProgress := fmt.Sprintf("%d/%d", batchNumber, totalBatch)

		batchLogger := applog.WithBatch(
			logger,
			total,
			batchProgress,
			len(batch),
			helper.Sample(batch, constant.SampleDatabaseSize),
		)

		if err := s.airlineRepo.Upsert(ctx, &batch, options); err != nil {
			batchLogger.Errorw(
				"Fail to upsert airline batch",
				"first_iata_code", batch[0].IATACode,
				"last_iata_code", batch[len(batch)-1].IATACode,
				"error", err,
			)

			return fmt.Errorf("upsert airline batch %s: %w", batchProgress, err)
		}

		batchLogger.Info("Upsert airline batch successfully")
	}

	return nil
}
