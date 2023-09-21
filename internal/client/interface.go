package client

import "context"

type BaseClient interface {
	DoGetURIWithQuery(ctx context.Context, url string, body map[string]string) (int64, []byte, error)
	makeEmptyBodyRequest(ctx context.Context, method string, url string) (int64, []byte, error)
}
