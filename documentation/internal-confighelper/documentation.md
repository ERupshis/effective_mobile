package confighelper // import "github.com/erupshis/effective_mobile/internal/confighelper"


FUNCTIONS

func Atoi64(value string) (int64, error)
    Atoi64 converts string into int64.

func SetEnvToParamIfNeed(param interface{}, val string) error
    SetEnvToParamIfNeed parses and assign string val to pointer of 'param'
    hidden under interface type. In case of unknown type returns error.

