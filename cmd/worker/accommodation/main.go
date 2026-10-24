package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/bootstrap/worker"
	"travel-crawler/internal/config"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/database"
	"travel-crawler/internal/kafka"
	"travel-crawler/internal/redis"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := applog.New(cfg.Log, cfg.App.Name).Named(constant.LoggerWorker)

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

	kafkaCli, err := kafka.Connect(ctx, cfg.Kafka)
	if err != nil {
		logger.Errorw(
			"failed to connect kafka",
			"error", err,
		)
		return
	}

	defer func() {
		if err := kafkaCli.Close(); err != nil {
			logger.Warnw(
				"failed to close kafka",
				"error", err,
			)
		}
	}()

	app, err := worker.NewApp(db, cfg, cache, kafkaCli, logger)
	if err != nil {
		logger.Errorw(
			"failed to initialize worker",
			"error", err,
		)
		return
	}

	logger.Info("Accommodation worker started")

	if err := app.AccommodationWorker.Run(ctx); err != nil {
		logger.Errorw(
			"Accommodation worker stopped unexpectedly",
			"error", err,
		)
		return
	}

	logger.Info("Accommodation worker stopped")
}
