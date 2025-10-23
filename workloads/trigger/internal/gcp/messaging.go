package gcp

import (
	"context"

	"go.uber.org/zap"

	"github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/provider"
)

// GCPMessaging implements the QueueProvider interface for GCP
type GCPMessaging struct{}

// NewGCPMessaging creates a new GCP messaging provider
func NewGCPMessaging() *GCPMessaging {
	return &GCPMessaging{}
}

// TriggerByQueue handles Pub/Sub message processing
func (g *GCPMessaging) TriggerByQueue(ctx context.Context, logger *zap.Logger, config provider.QueueConfig) error {
	logger.Info("triggering GCP Pub/Sub processing",
		zap.String("provider", "gcp"),
		zap.String("project", config.ProjectID),
		zap.String("topic", config.Topic),
		zap.String("region", config.Region),
		zap.Bool("verbose", config.Verbose),
	)

	if config.Verbose {
		logger.Debug("verbose mode enabled for Pub/Sub processing")
	}

	return nil
}
