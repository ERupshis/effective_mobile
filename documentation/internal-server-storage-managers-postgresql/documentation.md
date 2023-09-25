package postgresql // import "github.com/erupshis/effective_mobile/internal/server/storage/managers/postgresql"

Package postgresql postgresql handling PostgreSQL database.

const SchemaName = "persons_data" ...
func CreatePostgreDB(ctx context.Context, cfg config.Config, queriesHandler QueriesHandler, ...) (managers.BaseStorageManager, error)
type QueriesHandler struct{ ... }
    func CreateHandler(log logger.BaseLogger) QueriesHandler
