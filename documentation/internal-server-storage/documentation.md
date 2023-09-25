package storage // import "github.com/erupshis/effective_mobile/internal/server/storage"

Package storage provides main storage interface and implementation.

type BaseStorage interface{ ... }
    func Create(manager managers.BaseStorageManager, log logger.BaseLogger) BaseStorage
