package helpers

import (
	"fmt"
	"strconv"

	"github.com/erupshis/effective_mobile/internal/logger"
)

func ExecuteWithLogError(callback func() error, log logger.BaseLogger) {
	if err := callback(); err != nil {
		log.Info("callback execution finished with error: %v", err)
	}
}

func InterfaceToString(i interface{}) (string, error) {
	switch v := i.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	default:
		return "", fmt.Errorf("unsupported type: %T", v)
	}
}
