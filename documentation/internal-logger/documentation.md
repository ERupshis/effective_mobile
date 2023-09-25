package logger // import "github.com/erupshis/effective_mobile/internal/logger"

Package loggerZap provides interface for loggerZap in project.

TYPES

type BaseLogger interface {
	// Sync flushing any buffered log entries.
	Sync()

	// Info generates 'info' level log.
	Info(msg string, fields ...interface{})

	// Printf interface for kafka's implementation.
	Printf(msg string, fields ...interface{})

	// LogHandler handler for requests logging.
	LogHandler(h http.Handler) http.Handler
}
    BaseLogger interface of used loggerZap.

func CreateZapLogger(level string) (BaseLogger, error)
    CreateZapLogger create method for zap logger.

