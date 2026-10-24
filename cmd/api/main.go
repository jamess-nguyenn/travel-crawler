package main

import (
	"log"
	"net/http"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/auth"
	"travel-crawler/internal/bootstrap/server"
	"travel-crawler/internal/config"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/database"
	"travel-crawler/internal/redis"
	"travel-crawler/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := applog.New(cfg.Log, cfg.App.Name).Named(constant.LoggerAPI)

	db, err := database.Connect(cfg.MySQL)
	if err != nil {
		logger.Errorw(
			"failed to connect database",
			"error", err,
		)
		return
	}

	defer func() {
		if err := database.Close(db); err != nil {
			logger.Warnw(
				"failed to close database",
				"error", err,
			)
		}
	}()

	jwtMiddleware, err := auth.New(cfg.Auth0)
	if err != nil {
		logger.Errorw(
			"failed to initialize jwt middleware",
			"error", err,
		)
		return
	}

	cache, err := redis.NewCache(cfg.Redis, cfg.App.Key)
	if err != nil {
		logger.Errorw(
			"failed to connect redis",
			"error", err,
		)
		return
	}

	defer func() {
		if err := cache.Close(); err != nil {
			logger.Warnw(
				"failed to close cache",
				"error", err,
			)
		}
	}()

	app := server.NewApp(db, cfg, cache, logger)

	handler := router.New(app, jwtMiddleware, logger)

	address := cfg.App.Address()

	logger.Infow(
		"Starting server successfully",
		"address", address,
	)

	if err := http.ListenAndServe(address, handler); err != nil {
		logger.Errorw(
			"failed to start server",
			"address", address,
			"error", err,
		)
	}
}
