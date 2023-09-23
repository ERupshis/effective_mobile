package qraphqlctrl

import (
	"github.com/graphql-go/graphql"
)

const (
	argId         = "id"
	argName       = "name"
	argSurname    = "surname"
	argPatronymic = "patronymic"
	argAge        = "age"
	argGender     = "gender"
	argCountry    = "country"
)

var commonConfigs = map[string]*graphql.ArgumentConfig{
	argId: {
		Type: graphql.Int,
	},
	argName: {
		Type: graphql.String,
	},
	argSurname: {
		Type: graphql.String,
	},
	argPatronymic: {
		Type: graphql.String,
	},
	argAge: {
		Type: graphql.Int,
	},
	argGender: {
		Type: graphql.String,
	},
	argCountry: {
		Type: graphql.String,
	},
}

func getFieldConfigArgument(fields []string) graphql.FieldConfigArgument {
	res := graphql.FieldConfigArgument{}

	for _, field := range fields {
		res[field] = commonConfigs[field]
	}

	return res
}
