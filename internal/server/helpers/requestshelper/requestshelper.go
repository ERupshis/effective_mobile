package requestshelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

var FieldsInPersonData = []string{"id", "name", "surname", "patronymic", "age", "gender", "country", "page_num", "page_size"}

func ParsePersonDataFromJSON(rawData []byte) (*datastructs.PersonData, error) {
	personData := &datastructs.PersonData{}
	if err := json.Unmarshal(rawData, &personData); err != nil {
		return nil, fmt.Errorf("parse data to JSON: %w", err)
	}

	return personData, nil
}

func ParseQueryValuesIntoMap(values url.Values) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	for _, fieldName := range FieldsInPersonData {
		if values.Has(fieldName) {
			res[fieldName] = values.Get(fieldName)
			values.Del(fieldName)
		}
	}

	var err error
	if len(res) == 0 && len(values) != 0 {
		err = fmt.Errorf("query contains only incorrect keys")
	}

	return res, err
}

func FilterValues(values map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for _, fieldName := range FieldsInPersonData {
		if _, ok := values[fieldName]; ok {
			res[fieldName] = values[fieldName]
		}
	}
	return res
}

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

func IsPersonDataValid(data *datastructs.PersonData, allFieldsToCheck bool) (bool, error) {
	errorCriticalMessage := ""
	errorNonCriticalMessage := ""

	if data.Name == "" {
		errorCriticalMessage += " name"
	}

	if data.Surname == "" {
		errorCriticalMessage += " surname"
	}

	if allFieldsToCheck {
		if data.Age <= 0 {
			errorNonCriticalMessage += " age"
		}

		if data.Gender == "" {
			errorNonCriticalMessage += " gender"
		}
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
