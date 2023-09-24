package manager

import (
	"context"
)

type BaseCacheManager interface {
	Add(ctx context.Context, key map[string]interface{}, val interface{}) error
	Has(ctx context.Context, key map[string]interface{}) (bool, error)
	Get(ctx context.Context, key map[string]interface{}) ([]byte, error)

	Close() error
	Flush(ctx context.Context) error
}
