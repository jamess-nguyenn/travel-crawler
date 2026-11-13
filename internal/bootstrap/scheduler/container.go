package scheduler

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"travel-crawler/internal/airline"
	"travel-crawler/internal/airport"
	"travel-crawler/internal/config"
	"travel-crawler/internal/database"
	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/redis"

	"travel-crawler/internal/kafka"
)

type Container struct {
	cfg *config.Config

	txManager *database.TxManager

	airlineRepo airline.AirlineRepository
	airportRepo airport.AirportRepository

	http   *httpclient.Client
	cache  *redis.Cache
	kafka  *kafka.Client
	logger *zap.SugaredLogger
}

func NewContainer(
	db *gorm.DB,
	cfg *config.Config,
	cache *redis.Cache,
	kafka *kafka.Client,
	logger *zap.SugaredLogger,
) *Container {
	return &Container{
		cfg: cfg,

		txManager: database.NewTxManager(db),

		airlineRepo: airline.NewRepository(db),
		airportRepo: airport.NewRepository(db),

		http:   httpclient.New(logger),
		cache:  cache,
		kafka:  kafka,
		logger: logger,
	}
}
