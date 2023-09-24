package helpers

import (
	"fmt"

	"github.com/erupshis/effective_mobile/internal/logger"
)

func ExecuteWithLogError(callback func() error, log logger.BaseLogger) {
	if err := callback(); err != nil {
		log.Info("callback execution finished with error: %v", err)
	}
}

func InterfaceToString(i interface{}) string {
	return fmt.Sprintf("%v", i)
}
