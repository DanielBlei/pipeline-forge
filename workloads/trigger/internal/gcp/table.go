package gcp

import (
	"context"

	"go.uber.org/zap"

	"github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/provider"
)

// GCPTable implements the TableProvider interface for GCP
type GCPTable struct{}

// NewGCPTable creates a new GCP table provider
func NewGCPTable() *GCPTable {
	return &GCPTable{}
}

// TriggerByTable handles BigQuery table event processing
func (g *GCPTable) TriggerByTable(ctx context.Context, logger *zap.Logger, config provider.TableConfig) error {
	logger.Info("triggering GCP BigQuery processing",
		zap.String("provider", "gcp"),
		zap.String("project", config.ProjectID),
		zap.String("dataset", config.Dataset),
		zap.String("table", config.Table),
		zap.String("region", config.Region),
		zap.Bool("verbose", config.Verbose),
	)

	if config.Verbose {
		logger.Debug("verbose mode enabled for BigQuery processing")
	}

	return nil
}
