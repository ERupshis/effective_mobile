package helpers // import "github.com/erupshis/effective_mobile/internal/helpers"

func ExecuteWithLogError(callback func() error, log logger.BaseLogger)
func InterfaceToString(i interface{}) string
