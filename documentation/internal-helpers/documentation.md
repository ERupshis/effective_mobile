package helpers // import "github.com/erupshis/effective_mobile/internal/helpers"


FUNCTIONS

func ExecuteWithLogError(callback func() error, log logger.BaseLogger)
    ExecuteWithLogError support method for defer functions call which should
    return error.

func InterfaceToString(i interface{}) string
    InterfaceToString simple converter any interface into string.

