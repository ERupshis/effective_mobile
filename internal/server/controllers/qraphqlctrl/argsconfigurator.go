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
	argPageNum    = "page_num"
	argPageSize   = "page_size"
)

// commonConfigs map of predefined GraphQL's arguments.
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
	argPageNum: {
		Type: graphql.Int,
	},
	argPageSize: {
		Type: graphql.Int,
	},
}

// getFieldConfigArgument provides GraphQL's arguments with predefined type by name.
func getFieldConfigArgument(fields []string) graphql.FieldConfigArgument {
	res := graphql.FieldConfigArgument{}

	for _, field := range fields {
		res[field] = commonConfigs[field]
	}

	return res
}
