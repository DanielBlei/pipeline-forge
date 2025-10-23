package provider

import (
	"context"

	"go.uber.org/zap"
)

// Provider represents supported cloud providers
type Provider string

const (
	ProviderGCP   Provider = "gcp"
	ProviderAWS   Provider = "aws"
	ProviderAzure Provider = "azure"
)

// String returns the string representation of the provider
func (p Provider) String() string {
	return string(p)
}

// StorageConfig holds configuration for bucket/storage operations
type StorageConfig struct {
	Provider  Provider
	ProjectID string
	Bucket    string
	Region    string
	Verbose   bool
}

// QueueConfig holds configuration for queue/messaging operations
type QueueConfig struct {
	Provider  Provider
	ProjectID string
	Topic     string
	Region    string
	Verbose   bool
}

// TableConfig holds configuration for table/analytics operations
type TableConfig struct {
	Provider  Provider
	ProjectID string
	Dataset   string
	Table     string
	Region    string
	Verbose   bool
}

// StorageProvider defines interface for storage/bucket operations
type StorageProvider interface {
	TriggerByBucket(ctx context.Context, logger *zap.Logger, config StorageConfig) error
}

// QueueProvider defines interface for queue/messaging operations
type QueueProvider interface {
	TriggerByQueue(ctx context.Context, logger *zap.Logger, config QueueConfig) error
}

// TableProvider defines interface for table/analytics operations
type TableProvider interface {
	TriggerByTable(ctx context.Context, logger *zap.Logger, config TableConfig) error
}

// ProviderConstructors holds constructor functions for a provider
type ProviderConstructors struct {
	Name    string
	Storage func() StorageProvider
	Queue   func() QueueProvider
	Table   func() TableProvider
}
