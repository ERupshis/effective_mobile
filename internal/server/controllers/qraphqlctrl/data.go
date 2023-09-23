package qraphqlctrl

import (
	"github.com/graphql-go/graphql"
)

func getPersonType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Person",
		Fields: graphql.Fields{
			"Id": &graphql.Field{
				Type: graphql.Int,
			},
			"Name": &graphql.Field{
				Type: graphql.String,
			},
			"Surname": &graphql.Field{
				Type: graphql.String,
			},
			"Patronymic": &graphql.Field{
				Type: graphql.String,
			},
			"Age": &graphql.Field{
				Type: graphql.Int,
			},
			"Gender": &graphql.Field{
				Type: graphql.String,
			},
			"Country": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
}
