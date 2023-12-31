package helpers

import (
	"fmt"

	"github.com/erupshis/effective_mobile/internal/logger"
)

// ExecuteWithLogError support method for defer functions call which should return error.
func ExecuteWithLogError(callback func() error, log logger.BaseLogger) {
	if err := callback(); err != nil {
		log.Info("callback execution finished with error: %v", err)
	}
}

// InterfaceToString simple converter any interface into string.
func InterfaceToString(i interface{}) string {
	return fmt.Sprintf("%v", i)
}
