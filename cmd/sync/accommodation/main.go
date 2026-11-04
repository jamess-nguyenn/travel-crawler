package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"travel-crawler/internal/applog"
	"travel-crawler/internal/bootstrap/scheduler"
	"travel-crawler/internal/config"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/kafka"
)

const (
	defaultLookbackMonth = 6
)

func main() {
	dateFrom, countryCodes, err := parseArgs()
	if err != nil {
		log.Printf("failed to parse arguments: %v", err)
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logger := applog.New(cfg.Log, cfg.App.Name).Named(constant.LoggerSync)

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

	app, err := scheduler.NewApp(kafkaCli, logger)
	if err != nil {
		logger.Errorw(
			"failed to initialize sync app",
			"error", err,
		)
		return
	}

	syncLogger := logger.With(
		"date_from", dateFrom,
		"countries", countryCodes,
	)

	syncLogger.Info("Accommodation sync started")

	app.AccommodationPublisher.Publish(ctx, &dateFrom, countryCodes)

	syncLogger.Info("Accommodation sync completed")
}

func parseArgs() (time.Time, []string, error) {
	if len(os.Args) > 3 {
		return time.Time{}, nil, usageError()
	}

	dateFrom := helper.NowUTC().AddDate(0, -defaultLookbackMonth, 0)

	if len(os.Args) >= 2 {
		value := os.Args[1]

		parsed, err := time.Parse(time.DateOnly, value)
		if err != nil {
			return time.Time{}, nil, fmt.Errorf("invalid dateFrom %q: expected YYYY-MM-DD", value)
		}

		dateFrom = parsed.UTC()
	}

	var countryCodes []string

	if len(os.Args) == 3 && strings.TrimSpace(os.Args[2]) != "" {
		for _, code := range strings.Split(os.Args[2], ",") {
			code = strings.ToUpper(strings.TrimSpace(code))

			if code == "" {
				continue
			}

			countryCodes = append(countryCodes, code)
		}
	}

	return dateFrom, countryCodes, nil
}

func usageError() error {
	return fmt.Errorf(
		"invalid accommodation sync arguments\n\n%s",
		usage(),
	)
}

func usage() string {
	return `Usage:
  go run ./cmd/sync/accommodation [dateFrom] [countries]

Arguments:
  dateFrom    Start date for accommodation sync
  countries   Optional comma-separated country codes

Formats:
  dateFrom    YYYY-MM-DD
  countries   VN,SG,TH

Defaults:
  dateFrom    6 months ago (UTC)
  countries   All configured countries

Examples:
  go run ./cmd/sync/accommodation
  go run ./cmd/sync/accommodation 2026-01-01
  go run ./cmd/sync/accommodation 2026-01-01 VN
  go run ./cmd/sync/accommodation 2026-01-01 VN,SG
  go run ./cmd/sync/accommodation 2026-01-01 VN,SG,TH
  go run ./cmd/sync/accommodation 2026-01-01 ""`
}
