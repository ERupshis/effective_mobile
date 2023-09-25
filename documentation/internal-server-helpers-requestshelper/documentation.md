package requestshelper // import "github.com/erupshis/effective_mobile/internal/server/helpers/requestshelper"

Package requestshelper provides functions for input data parsing, validation,
type conversion, modification.

var FieldsInPersonData = []string{ ... }
func FilterPageNumAndPageSize(values map[string]interface{}) (map[string]interface{}, int64, int64)
func FilterValues(values map[string]interface{}) map[string]interface{}
func IsPersonDataValid(data *datastructs.PersonData, allFieldsToCheck bool) (bool, error)
func ParsePersonDataFromJSON(rawData []byte) (*datastructs.PersonData, error)
func ParseQueryValuesIntoMap(values url.Values) (map[string]interface{}, error)
