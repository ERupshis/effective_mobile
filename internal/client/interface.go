package client

import "context"

//go:generate mockgen -destination=../../mocks/mock_BaseClient.go -package=mocks github.com/erupshis/effective_mobile/internal/client BaseClient
type BaseClient interface {
	DoGetURIWithQuery(ctx context.Context, url string, body map[string]string) (int64, []byte, error)
}
