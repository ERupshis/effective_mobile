package retryer // import "github.com/erupshis/effective_mobile/internal/retryer"

func RetryCallWithTimeout(ctx context.Context, log logger.BaseLogger, intervals []int, ...) (int64, []byte, error)
func RetryCallWithTimeoutErrorOnly(ctx context.Context, log logger.BaseLogger, intervals []int, ...) error
