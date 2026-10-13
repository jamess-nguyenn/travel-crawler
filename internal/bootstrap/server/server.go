package server

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"travel-crawler/internal/auth"
	"travel-crawler/internal/config"
	"travel-crawler/internal/user"
)

type App struct {
	AuthService auth.Service
	UserService user.Service
}

func NewApp(
	db *gorm.DB,
	cfg *config.Config,
	logger *zap.SugaredLogger,
) *App {
	app := &App{}

	ctr := NewContainer(
		db,
		cfg,
		logger,
	)

	modules := []Module{
		AuthModule{},
		UserModule{},
	}

	for _, module := range modules {
		module.Init(
			app,
			ctr,
		)
	}

	return app
}
