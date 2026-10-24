package scheduler

import (
	"go.uber.org/zap"

	"travel-crawler/internal/kafka"
)

type Container struct {
	kafka  *kafka.Client
	logger *zap.SugaredLogger
}

func NewContainer(
	kafka *kafka.Client,
	logger *zap.SugaredLogger,
) *Container {
	return &Container{
		kafka:  kafka,
		logger: logger,
	}
}
