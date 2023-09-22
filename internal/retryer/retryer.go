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
	var cancels []context.CancelFunc
	go func() {
		sleepTime := 0
		for _, t := range intervals {
			sleepTime += t
		}
		time.Sleep(time.Duration(sleepTime) * time.Second)

		for _, cancel := range cancels {
			cancel()
		}
	}()
	for _, interval := range intervals {
		ctxWithTime, cancel := context.WithTimeout(ctx, time.Duration(interval)*time.Second)
		cancels = append(cancels, cancel)
		status, body, err = callback(ctxWithTime)
		if err == nil {
			return status, body, nil
		}

		attempt++
		if log != nil {
			log.Info("attempt '%d' to postJSON failed with error: %v", attempt, err)
		}

		if !canRetryCall(err, repeatableErrors) {
			break
		}
	}

	return status, body, err
}

func RetryCallWithTimeoutErrorOnly(ctx context.Context, log logger.BaseLogger, intervals []int, repeatableErrors []error,
	callback func(context.Context) error) error {
	var err error

	if intervals == nil {
		intervals = defIntervals
	}

	attemptNum := 0
	var cancels []context.CancelFunc
	go func() {
		sleepTime := 0
		for _, t := range intervals {
			sleepTime += t
		}
		time.Sleep(time.Duration(sleepTime) * time.Second)

		for _, cancel := range cancels {
			cancel()
		}
	}()
	for _, interval := range intervals {
		ctxWithTime, cancel := context.WithTimeout(ctx, time.Duration(interval)*time.Second)
		cancels = append(cancels, cancel)
		err = callback(ctxWithTime)
		if err == nil {
			return nil
		}

		attemptNum++
		if log != nil {
			log.Info("attemptNum '%d' to postJSON failed with error: %v", attemptNum, err)
		}

		if !canRetryCall(err, repeatableErrors) {
			log.Info("this kind of error is not retriable: %w", err)
			break
		}
	}

	return err
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
