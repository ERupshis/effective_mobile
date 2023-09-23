package retryer

import (
	"context"
	"sync"
	"time"

	"github.com/erupshis/effective_mobile/internal/logger"
)

var defIntervals = []int{1, 3, 5}

type ctxStruct struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func RetryCallWithTimeout(ctx context.Context, log logger.BaseLogger, intervals []int, repeatableErrors []error,
	callback func(context.Context) (int64, []byte, error)) (int64, []byte, error) {
	var status int64
	var body []byte
	var err error

	if intervals == nil {
		intervals = defIntervals
	}

	attempt := 0
	var cancels []ctxStruct
	mu := sync.Mutex{}
	go func() {
		sleepTime := 0
		for _, t := range intervals {
			sleepTime += t
		}
		time.Sleep(time.Duration(sleepTime+2) * time.Second)

		mu.Lock()
		for _, ctxToCancel := range cancels {
			ctxToCancel.cancel()
		}
		mu.Unlock()
	}()
	for _, interval := range intervals {
		ctxWithTime, cancel := context.WithTimeout(ctx, time.Duration(interval)*time.Second)
		mu.Lock()
		cancels = append(cancels, ctxStruct{ctxWithTime, cancel})
		mu.Unlock()
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

	var cancels []ctxStruct
	mu := sync.Mutex{}
	go func() {
		sleepTime := 0
		for _, t := range intervals {
			sleepTime += t
		}
		time.Sleep(time.Duration(sleepTime+2) * time.Second)

		mu.Lock()
		for _, ctxToCancel := range cancels {
			ctxToCancel.cancel()
		}
		mu.Unlock()
	}()
	for _, interval := range intervals {
		ctxWithTime, cancel := context.WithTimeout(ctx, time.Duration(interval)*time.Second)
		mu.Lock()
		cancels = append(cancels, ctxStruct{ctxWithTime, cancel})
		mu.Unlock()
		err = callback(ctxWithTime)
		if err == nil {
			return nil
		}

		attemptNum++
		if log != nil {
			log.Info("attemptNum '%d' to postJSON failed with error: %v", attemptNum, err)
		}

		if !canRetryCall(err, repeatableErrors) {
			log.Info("this kind of error is not retriable: %v", err)
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
