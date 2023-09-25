package retryer // import "github.com/erupshis/effective_mobile/internal/retryer"

Package retryer implements repeat requests logic.

FUNCTIONS

func RetryCallWithTimeout(ctx context.Context, log logger.BaseLogger, intervals []int, repeatableErrors []error,
	callback func(context.Context) (int64, []byte, error)) (int64, []byte, error)
    RetryCallWithTimeout generates repeats of function call if error occurs.
    Args:
      - ctx(context.Context), log Logger.BaseLogger,
      - intervals([]int) - count of repeats and pause between them (secs.);
      - repeatableErrors([]error) - errors - reasons to make repeat call.
        If empty - any error is signal to repeat call;
      - callback(func(context.Context) (int64, []byte, error)) - function to
        call.

func RetryCallWithTimeoutErrorOnly(ctx context.Context, log logger.BaseLogger, intervals []int, repeatableErrors []error,
	callback func(context.Context) error) error
    RetryCallWithTimeoutErrorOnly generates repeats of function call if error
    occurs. Args:
      - ctx(context.Context), log Logger.BaseLogger,
      - intervals([]int) - count of repeats and pause between them (secs.);
      - repeatableErrors([]error) - errors - reasons to make repeat call.
        If empty - any error is signal to repeat call;
      - callback (func(context.Context) error) - function to call.

