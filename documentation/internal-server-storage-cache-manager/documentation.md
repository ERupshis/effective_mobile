package manager // import "github.com/erupshis/effective_mobile/internal/server/storage/cache/manager"

Package manager provides cache manager interface and implements it.

type BaseCacheManager interface{ ... }
    func CreateRedis(ctx context.Context, dsn string, log logger.BaseLogger) (BaseCacheManager, error)
