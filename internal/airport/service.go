package airport

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/data"
	"travel-crawler/internal/database/scope"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/model"
	"travel-crawler/internal/response"

	"travel-crawler/internal/database"
)

type Service interface {
	Get(
		ctx context.Context,
		req GetRequest,
	) (*Airport, error)

	List(
		ctx context.Context,
		req ListRequest,
	) (*response.ListResponse[Airport], error)
}

type SyncService interface {
	Sync(
		ctx context.Context,
	) error
}

type service struct {
	txManager   *database.TxManager
	airportRepo AirportRepository
	logger      *zap.SugaredLogger
}

type syncService struct {
	txManager   *database.TxManager
	airportRepo AirportRepository
	logger      *zap.SugaredLogger
}

func NewService(
	txManager *database.TxManager,
	airportRepo AirportRepository,
	logger *zap.SugaredLogger,
) Service {
	return &service{
		txManager:   txManager,
		airportRepo: airportRepo,
		logger:      logger,
	}
}

func NewSyncService(
	txManager *database.TxManager,
	airportRepo AirportRepository,
	logger *zap.SugaredLogger,
) SyncService {
	return &syncService{
		txManager:   txManager,
		airportRepo: airportRepo,
		logger:      logger,
	}
}

func (s *service) Get(
	ctx context.Context,
	req GetRequest,
) (*Airport, error) {
	logger := applog.FromContext(s.logger, ctx)

	conditions := map[string]any{"id": req.ID}

	airport, err := s.airportRepo.Get(ctx, conditions, selectColumns)

	getLogger := applog.WithGet(
		logger,
		conditions,
		selectColumns,
		airport,
	)

	if err != nil {
		getLogger.Errorw(
			"Fail to get airport",
			"error", err,
		)

		return nil, err
	}

	getLogger.Info("Get airport successfully")

	return airport, nil
}

func (s *service) buildListQuery(
	req ListRequest,
) database.ListQuery {
	conditions := make(map[string]any)

	filters := []scope.Scope{
		scope.SearchColumns([]string{"iata_code", "city", "country"}, req.Search),
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
) (*response.ListResponse[Airport], error) {
	logger := applog.FromContext(s.logger, ctx)

	listQuery := s.buildListQuery(req)

	airports, total, err := s.airportRepo.List(ctx, listQuery)

	listLogger := applog.WithList(
		logger,
		listQuery,
		total,
		airports,
	)

	if err != nil {
		listLogger.Errorw(
			"Fail to list airports",
			"error", err,
		)

		return nil, err
	}

	listLogger.Info("List airports successfully")

	return &response.ListResponse[Airport]{
		Items:      airports,
		Pagination: response.NewPagination(req.GetPage(), req.GetPageSize(), total),
		Request:    req,
	}, nil
}

func (s *syncService) Sync(
	ctx context.Context,
) error {
	logger := applog.FromContext(s.logger, ctx)

	airports, err := data.Airports()
	if err != nil {
		return fmt.Errorf("get airports: %w", err)
	}

	total := len(airports)

	listLogger := logger.With(
		"total", total,
		"data", helper.Sample(airports, constant.SampleLogSize),
	)

	listLogger.Info("Search airports successfully")

	if total == 0 {
		listLogger.Info("No airports to sync")
		return nil
	}

	if err := s.persistAirport(ctx, mapAirports(airports), logger); err != nil {
		return fmt.Errorf(
			"persist airports from %s: %w",
			err,
		)
	}

	listLogger.Info("Sync airports successfully")

	return nil
}

func mapAirport(item data.Airport, now time.Time) Airport {
	return Airport{
		Model: model.Model{
			Metadata:  helper.ToJSONObject(item),
			CreatedAt: now,
			CreatedBy: constant.SystemUserID,
			UpdatedAt: now,
			UpdatedBy: constant.SystemUserID,
		},
		IATACode: item.IATA,
		City:     item.City,
		Country:  item.Country,
	}
}

func mapAirports(items []data.Airport) []Airport {
	now := helper.NowUTC()

	return helper.MapSlice(items, func(item data.Airport) Airport {
		return mapAirport(item, now)
	})
}

func (s *syncService) persistAirport(
	ctx context.Context,
	airports []Airport,
	logger *zap.SugaredLogger,
) error {
	total := len(airports)
	if total == 0 {
		return nil
	}

	options := database.UpsertOptions{
		ConflictColumns: upsertConflictColumns,
		UpdateColumns:   upsertUpdateColumns,
	}

	batches := helper.Chunk(airports, constant.DatabaseBatchSize)
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

		if err := s.airportRepo.Upsert(ctx, &batch, options); err != nil {
			batchLogger.Errorw(
				"Fail to upsert airport batch",
				"first_iata_code", batch[0].IATACode,
				"last_iata_code", batch[len(batch)-1].IATACode,
				"error", err,
			)

			return fmt.Errorf("upsert airport batch %s: %w", batchProgress, err)
		}

		batchLogger.Info("Upsert airport batch successfully")
	}

	return nil
}
