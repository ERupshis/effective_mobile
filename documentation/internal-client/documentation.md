package client // import "github.com/erupshis/effective_mobile/internal/client"


TYPES

type BaseClient interface {
	DoGetURIWithQuery(ctx context.Context, url string, body map[string]string) (int64, []byte, error)
}

func CreateDefault(log logger.BaseLogger) BaseClient

type DefaultClient struct {
	// Has unexported fields.
}

func (c *DefaultClient) DoGetURIWithQuery(ctx context.Context, url string, params map[string]string) (int64, []byte, error)

