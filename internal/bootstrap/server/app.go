package server

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"travel-crawler/internal/accommodation"
	"travel-crawler/internal/airline"
	"travel-crawler/internal/airport"
	"travel-crawler/internal/auth"
	"travel-crawler/internal/config"
	"travel-crawler/internal/redis"
	"travel-crawler/internal/user"
)

type App struct {
	AuthService          auth.Service
	UserService          user.Service
	AccommodationService accommodation.Service
	AirlineService       airline.Service
	AirportService       airport.Service
}

func NewApp(
	db *gorm.DB,
	cfg *config.Config,
	cache *redis.Cache,
	logger *zap.SugaredLogger,
) *App {
	app := &App{}

	ctr := NewContainer(db, cfg, cache, logger)

	modules := []Module{
		AuthModule{},
		UserModule{},
		AccommodationModule{},
		AirlineModule{},
		AirportModule{},
	}

	for _, module := range modules {
		module.Init(app, ctr)
	}

	return app
}
