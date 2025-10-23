package provider

import (
	"fmt"
	"sync"
)

// Registry holds registered providers
type Registry struct {
	mu        sync.RWMutex
	providers map[Provider]ProviderConstructors
}

// DefaultRegistry is the global provider registry
var DefaultRegistry = &Registry{
	providers: make(map[Provider]ProviderConstructors),
}

// Register adds a provider to the registry
func (r *Registry) Register(p Provider, constructors ProviderConstructors) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p] = constructors
}

// GetStorageProvider returns a storage provider for the given provider type
func (r *Registry) GetStorageProvider(p Provider) (StorageProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	constructors, ok := r.providers[p]
	if !ok {
		return nil, fmt.Errorf("unsupported storage provider: %s (supported: %v)", p, r.SupportedProviders())
	}

	if constructors.Storage == nil {
		return nil, fmt.Errorf("storage not implemented for provider: %s", p)
	}

	return constructors.Storage(), nil
}

// GetQueueProvider returns a queue provider for the given provider type
func (r *Registry) GetQueueProvider(p Provider) (QueueProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	constructors, ok := r.providers[p]
	if !ok {
		return nil, fmt.Errorf("unsupported queue provider: %s (supported: %v)", p, r.SupportedProviders())
	}

	if constructors.Queue == nil {
		return nil, fmt.Errorf("queue not implemented for provider: %s", p)
	}

	return constructors.Queue(), nil
}

// GetTableProvider returns a table provider for the given provider type
func (r *Registry) GetTableProvider(p Provider) (TableProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	constructors, ok := r.providers[p]
	if !ok {
		return nil, fmt.Errorf("unsupported table provider: %s (supported: %v)", p, r.SupportedProviders())
	}

	if constructors.Table == nil {
		return nil, fmt.Errorf("table not implemented for provider: %s", p)
	}

	return constructors.Table(), nil
}

// SupportedProviders returns a list of supported providers
func (r *Registry) SupportedProviders() []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	providers := make([]Provider, 0, len(r.providers))
	for p := range r.providers {
		providers = append(providers, p)
	}
	return providers
}

// Convenience functions for default registry
func GetStorageProvider(p Provider) (StorageProvider, error) {
	return DefaultRegistry.GetStorageProvider(p)
}

func GetQueueProvider(p Provider) (QueueProvider, error) {
	return DefaultRegistry.GetQueueProvider(p)
}

func GetTableProvider(p Provider) (TableProvider, error) {
	return DefaultRegistry.GetTableProvider(p)
}
