package qraphqlctrl

import (
	"fmt"
	"net/http"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

const packageName = "graphqlctrl"

var personType = getPersonType()

type Controller struct {
	persons []datastructs.PersonData

	strg storage.BaseStorageManager
	log  logger.BaseLogger
}

func Create(strg storage.BaseStorageManager, log logger.BaseLogger) *Controller {
	return &Controller{
		strg: strg,
		log:  log,
	}
}

func (c *Controller) Route() *chi.Mux {
	r := chi.NewRouter()
	graphHandler, err := c.createHandler()
	if err != nil {
		c.log.Info("["+packageName+":Controller:Route] failed to create route: %v", err)
		r.HandleFunc("/", c.InternalErrorHandler)
		return r
	}

	r.Handle("/", graphHandler)
	return r
}

func (c *Controller) InternalErrorHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "graphql currently unavailable", http.StatusInternalServerError)
}

func (c *Controller) createHandler() (*handler.Handler, error) {
	schema, err := c.createSchema()
	if err != nil {
		return nil, fmt.Errorf("create schema: %w", err)
	}

	return handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	}), nil
}

func (c *Controller) createSchema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    c.createQueries(),
		Mutation: c.createMutations(),
	})
}

func (c *Controller) createQueries() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"persons": c.readPersonsQuery(),
		},
	})
}

func (c *Controller) createMutations() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createPerson": c.createPersonMutation(),
			"updatePerson": c.updatePersonMutation(),
			"deletePerson": c.deletePersonMutation(),
		},
	})
}

func (c *Controller) readPersonsQuery() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(personType),
		Args: getFieldConfigArgument([]string{
			argName,
			argSurname,
			argAge,
			argGender,
		}),
		Resolve: c.readPersonsResolver,
	}
}

func (c *Controller) createPersonMutation() *graphql.Field {
	return &graphql.Field{
		Type: personType,
		Args: getFieldConfigArgument([]string{
			argName,
			argSurname,
			argPatronymic,
			argAge,
			argGender,
			argCountry,
		}),
		Resolve: c.createPersonResolver,
	}
}

func (c *Controller) updatePersonMutation() *graphql.Field {
	return &graphql.Field{
		Type: personType,
		Args: getFieldConfigArgument([]string{
			argId,
			argName,
			argSurname,
			argPatronymic,
			argAge,
			argGender,
			argCountry,
		}),
		Resolve: c.updatePersonResolver,
	}
}

func (c *Controller) deletePersonMutation() *graphql.Field {
	return &graphql.Field{
		Type: personType,
		Args: getFieldConfigArgument([]string{
			argId,
		}),
		Resolve: c.deletePersonResolver,
	}
}
