package requestshelper // import "github.com/erupshis/effective_mobile/internal/server/helpers/requestshelper"

Package requestshelper provides functions for input data parsing, validation,
type conversion, modification.

VARIABLES

var FieldsInPersonData = []string{"id", "name", "surname", "patronymic", "age", "gender", "country", "page_num", "page_size"}
    FieldsInPersonData slice of available values in query.


FUNCTIONS

func FilterPageNumAndPageSize(values map[string]interface{}) (map[string]interface{}, int64, int64)
    FilterPageNumAndPageSize extracts pageNum and pageSize values from values
    map, otherwise returns 0, 0.

func FilterValues(values map[string]interface{}) map[string]interface{}
    FilterValues filters maps[string]interface{} values in and leave values
    with key equal to some from FieldsInPersonData. Also lowers register of all
    string values.

func IsPersonDataValid(data *datastructs.PersonData, allFieldsToCheck bool) (bool, error)
    IsPersonDataValid validates person data structs and check fields or
    non-empty value. Patronymic is ignored.

func ParsePersonDataFromJSON(rawData []byte) (*datastructs.PersonData, error)
    ParsePersonDataFromJSON parse request body(json type) into person data
    struct.

func ParseQueryValuesIntoMap(values url.Values) (map[string]interface{}, error)
    ParseQueryValuesIntoMap parse url.Values int map, filters and leave only
    queries with keys from FieldsInPersonData, Generates error in case of
    incorrect queries in request.

