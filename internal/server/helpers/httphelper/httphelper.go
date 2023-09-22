package httphelper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

var FieldsInPersonData = []string{"name", "surname", "patronymic", "age", "gender", "country"}

func ParsePersonDataFromJSON(rawData []byte) (*datastructs.PersonData, error) {
	personData := &datastructs.PersonData{}
	if err := json.Unmarshal(rawData, &personData); err != nil {
		return nil, fmt.Errorf("parse data to JSON: %w", err)
	}

	return personData, nil
}

func ParseQueryValuesIntoMap(values url.Values) (map[string]string, error) {
	res := map[string]string{}
	for _, fieldName := range FieldsInPersonData {
		if values.Has(fieldName) {
			res[fieldName] = values.Get(fieldName)
		}
	}

	var err error
	if len(res) == 0 {
		err = fmt.Errorf("missing any siutable query value")
	}

	return res, err
}

func ParsePageAndPageSize(values url.Values) (int64, int64) {
	rawPage := values.Get("page")
	rawPageSize := values.Get("pageSize")

	page, _ := strconv.Atoi(rawPage)
	pageSize, _ := strconv.Atoi(rawPageSize)

	return int64(page), int64(pageSize)
}

func IsPersonDataValid(data *datastructs.PersonData, allFieldsToCheck bool) (bool, error) {
	errorCriticalMessage := ""
	errorNonCriticalMessage := ""

	if data.Name == "" {
		errorCriticalMessage += " name"
	}

	if data.Name == "" {
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
