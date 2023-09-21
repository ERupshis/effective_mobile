package httphelper

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

var fieldsInPersonData = []string{"name", "surname", "age", "patronymic", "gender", "nationality"}

func ParseQueryValues(values url.Values) (*datastructs.PersonData, error) {
	personData := &datastructs.PersonData{
		Name:        values.Get("name"),
		Surname:     values.Get("surname"),
		Patronymic:  values.Get("patronymic"),
		Gender:      values.Get("gender"),
		Nationality: values.Get("nationality"),
	}

	var err error
	personData.Age, err = strconv.ParseInt(values.Get("age"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse 'age' in query: %w", err)
	}

	return personData, nil
}

func ParseQueryValuesIntoMap(values url.Values) (map[string]string, error) {
	res := map[string]string{}
	for _, fieldName := range fieldsInPersonData {
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
		if data.Nationality == "" {
			errorNonCriticalMessage += " nationality"
		}
	}

	allFieldsEmpty := data.Name == "" && data.Surname == "" && data.Gender == "" && data.Nationality == "" && data.Age <= 0
	if allFieldsEmpty {
		return false, fmt.Errorf("all person data fields empty")
	}

	if allFieldsToCheck && (errorCriticalMessage != "" || errorNonCriticalMessage != "") {
		return false, fmt.Errorf("person data is not valid: %s", errorCriticalMessage+errorNonCriticalMessage)
	}

	return true, nil
}
