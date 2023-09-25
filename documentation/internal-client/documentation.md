package client // import "github.com/erupshis/effective_mobile/internal/client"

type BaseClient interface{ ... }
    func CreateDefault(log logger.BaseLogger) BaseClient
type DefaultClient struct{ ... }
