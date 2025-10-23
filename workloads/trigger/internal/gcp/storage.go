package gcp

import (
	"context"

	"go.uber.org/zap"

	"github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/provider"
)

// GCPStorage implements the StorageProvider interface for GCP
type GCPStorage struct{}

// NewGCPStorage creates a new GCP storage provider
func NewGCPStorage() *GCPStorage {
	return &GCPStorage{}
}

// TriggerByBucket handles GCS bucket event processing
func (g *GCPStorage) TriggerByBucket(ctx context.Context, logger *zap.Logger, config provider.StorageConfig) error {
	logger.Info("triggering GCP bucket processing",
		zap.String("provider", "gcp"),
		zap.String("project", config.ProjectID),
		zap.String("bucket", config.Bucket),
		zap.String("region", config.Region),
		zap.Bool("verbose", config.Verbose),
	)

	if config.Verbose {
		logger.Debug("verbose mode enabled for GCS bucket processing")
	}

	return nil
}
