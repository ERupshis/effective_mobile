// Package requestshelper provides functions for input data parsing, validation, type conversion, modification.
package requestshelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

// FieldsInPersonData slice of available values in query.
var FieldsInPersonData = []string{"id", "name", "surname", "patronymic", "age", "gender", "country", "page_num", "page_size"}

// ParsePersonDataFromJSON parse request body(json type) into person data struct.
func ParsePersonDataFromJSON(rawData []byte) (*datastructs.PersonData, error) {
	personData := &datastructs.PersonData{}
	if err := json.Unmarshal(rawData, &personData); err != nil {
		return nil, fmt.Errorf("parse data to JSON: %w", err)
	}

	return personData, nil
}

// ParseQueryValuesIntoMap parse url.Values int map, filters and leave only queries with keys from FieldsInPersonData,
// Generates error in case of incorrect queries in request.
func ParseQueryValuesIntoMap(values url.Values) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	for _, fieldName := range FieldsInPersonData {
		if values.Has(fieldName) {
			res[fieldName] = strings.ToLower(values.Get(fieldName))
			values.Del(fieldName)
		}
	}

	var err error
	if len(res) == 0 && len(values) != 0 {
		err = fmt.Errorf("query contains only incorrect keys")
	}

	return res, err
}

// IsPersonDataValid validates person data structs and check fields or non-empty value.
// Patronymic is ignored.
func IsPersonDataValid(data *datastructs.PersonData, allFieldsToCheck bool) (bool, error) {
	errorCriticalMessage := ""
	errorNonCriticalMessage := ""

	data.Name = strings.ToLower(data.Name)
	if data.Name == "" {
		errorCriticalMessage += " name"
	}

	data.Surname = strings.ToLower(data.Surname)
	if data.Surname == "" {
		errorCriticalMessage += " surname"
	}

	data.Patronymic = strings.ToLower(data.Patronymic)

	if allFieldsToCheck {
		if data.Age <= 0 {
			errorNonCriticalMessage += " age"
		}

		data.Gender = strings.ToLower(data.Gender)
		if data.Gender == "" {
			errorNonCriticalMessage += " gender"
		}

		data.Country = strings.ToLower(data.Country)
		if data.Country == "" {
			errorNonCriticalMessage += " country"
		}
	}

	allFieldsEmpty := data.Name == "" && data.Surname == "" && data.Gender == "" && data.Country == "" && data.Age <= 0
	if allFieldsEmpty {
		return false, fmt.Errorf("all person data fields empty")
	}

	if allFieldsToCheck && (errorCriticalMessage != "" || errorNonCriticalMessage != "") {
		return false, fmt.Errorf("person data is not valid: %s", errorCriticalMessage+errorNonCriticalMessage)
	}

	return true, nil
}

// FilterValues filters maps[string]interface{} values in and leave values with key equal to some from FieldsInPersonData.
// Also lowers register of all string values.
func FilterValues(values map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for _, fieldName := range FieldsInPersonData {
		if _, ok := values[fieldName]; ok {
			res[fieldName] = strings.ToLower(ConvertQueryValueIntoString(values[fieldName]))
		}
	}
	return res
}

// FilterPageNumAndPageSize extracts pageNum and pageSize values from values map, otherwise returns 0, 0.
func FilterPageNumAndPageSize(values map[string]interface{}) (map[string]interface{}, int64, int64) {
	var pageNum int64
	var pageSize int64

	if val, ok := values["page_num"]; ok {
		pageNum = convertQueryValueIntoInt64(val)
		delete(values, "page_num")
	}

	if val, ok := values["page_size"]; ok {
		pageSize = convertQueryValueIntoInt64(val)
		delete(values, "page_size")
	}

	return values, pageNum, pageSize
}

// convertQueryValueIntoInt64 converts any value into int64, if possible. Otherwise returns 0.
func convertQueryValueIntoInt64(value interface{}) int64 {
	var res int64
	switch param := value.(type) {
	case string:
		intVal, err := strconv.Atoi(param)
		if err != nil {
			res = 0
		} else {
			res = int64(intVal)
		}
	case int:
		res = int64(param)
	default:
		res = 0
	}

	return res
}

// ConvertQueryValueIntoString provides some value into string representation (for interface{} types).
func ConvertQueryValueIntoString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
