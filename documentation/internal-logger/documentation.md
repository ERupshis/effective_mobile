package logger // import "github.com/erupshis/effective_mobile/internal/logger"

type BaseLogger interface{ ... }
    func CreateZapLogger(level string) (BaseLogger, error)
type Logger struct{ ... }
