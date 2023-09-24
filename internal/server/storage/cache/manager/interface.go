// Package manager provides cache manager interface and implements it.
package manager

import (
	"context"
)

// BaseCacheManager cache manager interface.
type BaseCacheManager interface {
	// Add adds new value in cache.
	Add(ctx context.Context, key map[string]interface{}, val interface{}) error
	// Has checks that provided keys was added in cache.
	Has(ctx context.Context, key map[string]interface{}) (bool, error)
	// Get returns value for key from cache if exists.
	Get(ctx context.Context, key map[string]interface{}) ([]byte, error)

	// Close closes connection to cache.
	Close() error

	// Flush resets cache.
	Flush(ctx context.Context) error
}
