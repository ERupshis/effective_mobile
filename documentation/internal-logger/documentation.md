package logger // import "github.com/erupshis/effective_mobile/internal/logger"


TYPES

type BaseLogger interface {
	Sync()

	Info(msg string, fields ...interface{})
	Printf(msg string, fields ...interface{})
	LogHandler(h http.Handler) http.Handler
}

func CreateZapLogger(level string) (BaseLogger, error)

type Logger struct {
	// Has unexported fields.
}

func (l *Logger) Info(msg string, fields ...interface{})

func (l *Logger) LogHandler(h http.Handler) http.Handler

func (l *Logger) Printf(msg string, fields ...interface{})

func (l *Logger) Sync()

