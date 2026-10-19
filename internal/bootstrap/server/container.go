package server

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"travel-crawler/internal/accommodation"
	"travel-crawler/internal/auth0"
	"travel-crawler/internal/config"
	"travel-crawler/internal/database"
	"travel-crawler/internal/httpclient"
	"travel-crawler/internal/redis"
	"travel-crawler/internal/user"
)

type Container struct {
	txManager *database.TxManager

	userRepo          user.UserRepository
	accommodationRepo accommodation.AccommodationRepository

	auth0Client auth0.Client

	logger *zap.SugaredLogger
}

func NewContainer(
	db *gorm.DB,
	cfg *config.Config,
	cache *redis.Cache,
	logger *zap.SugaredLogger,
) *Container {
	http := httpclient.New(logger)

	return &Container{
		txManager: database.NewTxManager(db),

		userRepo:          user.NewRepository(db),
		accommodationRepo: accommodation.NewRepository(db),

		auth0Client: auth0.New(http, cfg.Auth0, cache, logger),

		logger: logger,
	}
}
