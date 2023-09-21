package retryer

import (
	"context"
	"time"

	"github.com/erupshis/effective_mobile/internal/logger"
)

var defIntervals = []int{1, 3, 5}

func RetryCallWithTimeout(ctx context.Context, log logger.BaseLogger, intervals []int, repeatableErrors []error,
	callback func(context.Context) (int64, []byte, error)) (int64, []byte, error) {
	var status int64
	var body []byte
	var err error

	if intervals == nil {
		intervals = defIntervals
	}

	attempt := 0
	for _, interval := range intervals {
		ctxWithTime, cancel := context.WithTimeout(ctx, time.Duration(interval)*time.Second)
		status, body, err = callback(ctxWithTime)
		if err == nil {
			cancel()
			return status, body, nil
		}

		<-ctxWithTime.Done()

		attempt++
		if log != nil {
			log.Info("attempt '%d' to postJSON failed with error: %v", attempt, err)
		}

		if !canRetryCall(err, repeatableErrors) {
			cancel()
			break
		}
		cancel()
	}

	return status, body, err
}

func canRetryCall(err error, repeatableErrors []error) bool {
	if repeatableErrors == nil {
		return true
	}

	canRetry := false
	for _, repeatableError := range repeatableErrors {
		if err.Error() == repeatableError.Error() {
			canRetry = true
		}
	}

	return canRetry
}
