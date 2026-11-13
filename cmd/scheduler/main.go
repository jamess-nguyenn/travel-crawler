package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/bootstrap/scheduler"
	"travel-crawler/internal/config"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/database"
	"travel-crawler/internal/helper"
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

	logger := applog.New(cfg.Log, cfg.App.Name).Named(constant.LoggerScheduler)

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

	app, err := scheduler.NewApp(db, cfg, cache, kafkaCli, logger)
	if err != nil {
		logger.Errorw(
			"failed to initialize cron scheduler",
			"error", err,
		)
		return
	}

	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)

	if _, err := c.AddFunc("0 */10 * * * *", func() {
		dateFrom := helper.NowUTC().Add(-30 * time.Minute)

		app.AccommodationPublisher.Publish(ctx, &dateFrom, nil)
	}); err != nil {
		logger.Errorw(
			"failed to register accommodation crawl schedule",
			"error", err,
		)
		return
	}

	if _, err := c.AddFunc("0 0 1 * * *", func() {
		app.AirlineSyncJob.Run(ctx)
	}); err != nil {
		logger.Errorw(
			"failed to register airline crawl schedule",
			"error", err,
		)
		return
	}

	c.Start()

	logger.Info("Cron scheduler started")

	<-ctx.Done()

	logger.Info("Stopping cron scheduler")

	stopCtx := c.Stop()
	<-stopCtx.Done()

	logger.Info("Cron scheduler stopped")
}
