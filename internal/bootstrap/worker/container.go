package worker

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"travel-crawler/internal/accommodation"
	"travel-crawler/internal/config"
	"travel-crawler/internal/database"
	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/kafka"
	"travel-crawler/internal/redis"
)

type Container struct {
	cfg *config.Config

	txManager *database.TxManager

	accommodationRepo accommodation.AccommodationRepository

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

		accommodationRepo: accommodation.NewRepository(db),

		http:   httpclient.New(logger),
		cache:  cache,
		kafka:  kafka,
		logger: logger,
	}
}
