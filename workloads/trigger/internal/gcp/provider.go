package gcp

import (
	"github.com/DanielBlei/pipeline-forge/workloads/trigger/internal/provider"
)

// init registers the GCP provider with the global registry
func init() {
	provider.DefaultRegistry.Register(provider.ProviderGCP, provider.ProviderConstructors{
		Name:    "Google Cloud Platform",
		Storage: func() provider.StorageProvider { return NewGCPStorage() },
		Queue:   func() provider.QueueProvider { return NewGCPMessaging() },
		Table:   func() provider.TableProvider { return NewGCPTable() },
	})
}
