package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"travel-crawler/internal/applog"
	"travel-crawler/internal/constant"
	"travel-crawler/internal/helper"
	"travel-crawler/internal/requestcontext"

	"travel-crawler/internal/config"
	"travel-crawler/internal/database"
)

const (
	commandUp     = "up"
	commandDown   = "down"
	commandReset  = "reset"
	commandStatus = "status"

	defaultDownStep = 1
)

func main() {
	if err := run(); err != nil {
		log.Printf("migration command failed: %v", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return usageError()
	}

	ctx := requestcontext.WithNewRequestContext(
		context.Background(),
	)

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger := applog.FromContext(
		applog.New(cfg.Log, cfg.App.Name), ctx,
	).Named(constant.LoggerMigration)

	db, err := database.Connect(cfg.MySQL)
	if err != nil {
		logger.Errorw(
			"failed to connect database",
			"error", err,
		)

		return fmt.Errorf("failed to connect database: %w", err)
	}

	defer func() {
		if err := database.Close(db); err != nil {
			logger.Warnw(
				"failed to close database",
				"error", err,
			)
		}
	}()

	command := os.Args[1]

	logger.Infow(
		"starting migration command",
		"command", command,
	)

	switch command {
	case commandUp:
		logger.Info("applying pending migrations")

		if err := database.Up(db); err != nil {
			logger.Errorw(
				"failed to apply migrations",
				"error", err,
			)

			return err
		}

		logger.Info("migrations applied successfully")

		return nil

	case commandDown:
		step, err := parseDownStep()
		if err != nil {
			logger.Errorw(
				"failed to parse migration down step",
				"error", err,
			)

			return err
		}

		logger.Infow(
			"rolling back migrations",
			"step", step,
		)

		if err := database.Down(db, step); err != nil {
			logger.Errorw(
				"failed to roll back migrations",
				"step", step,
				"error", err,
			)

			return err
		}

		logger.Infow(
			"migrations rolled back successfully",
			"step", step,
		)

		return nil

	case commandReset:
		logger.Info("resetting migrations")

		if err := database.Reset(db); err != nil {
			logger.Errorw(
				"failed to reset migrations",
				"error", err,
			)

			return err
		}

		logger.Info("migrations reset successfully")

		return nil

	case commandStatus:
		logger.Info("checking migration status")

		if err := database.Status(db); err != nil {
			logger.Errorw(
				"failed to check migration status",
				"error", err,
			)

			return err
		}

		logger.Info("migration status completed")

		return nil

	default:
		logger.Errorw(
			"unknown migration command",
			"command", command,
		)

		return fmt.Errorf(
			"unknown migration command %q\n\n%s",
			command,
			usage(),
		)

	}
}

func parseDownStep() (int, error) {
	if len(os.Args) < 3 {
		return defaultDownStep, nil
	}

	step, err := helper.ParseInt(os.Args[2])
	if err != nil {
		return 0, fmt.Errorf(
			"invalid migration down step %q: must be an integer",
			os.Args[2],
		)
	}

	if step < 1 {
		return 0, fmt.Errorf(
			"migration down step must be greater than 0",
		)
	}

	return step, nil
}

func usageError() error {
	return fmt.Errorf(
		"missing migration command\n\n%s",
		usage(),
	)
}

func usage() string {
	return `Usage:
  go run ./cmd/migration <command> [step]

Commands:
  up             Apply all pending migrations
  down [step]    Roll back migrations
                 Default step: 1
  reset          Roll back all migrations
  status         Show migration status

Examples:
  go run ./cmd/migration up
  go run ./cmd/migration down
  go run ./cmd/migration down 3
  go run ./cmd/migration reset
  go run ./cmd/migration status`
}
