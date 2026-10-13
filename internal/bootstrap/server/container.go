package server

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"travel-crawler/internal/auth0"
	"travel-crawler/internal/config"
	"travel-crawler/internal/database"
	"travel-crawler/internal/user"
)

type Container struct {
	txManager *database.TxManager

	userRepo user.UserRepository

	auth0Client auth0.Client

	logger *zap.SugaredLogger
}

func NewContainer(
	db *gorm.DB,
	cfg *config.Config,
	logger *zap.SugaredLogger,
) *Container {
	return &Container{
		txManager: database.NewTxManager(db),

		userRepo: user.NewRepository(db),

		auth0Client: auth0.New(cfg.Auth0, logger),

		logger: logger,
	}
}
