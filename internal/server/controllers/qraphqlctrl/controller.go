// Package qraphqlctrl provides GraphQL handling.
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
	strg storage.BaseStorage
	log  logger.BaseLogger
}

// Create returns controller.
func Create(strg storage.BaseStorage, log logger.BaseLogger) *Controller {
	return &Controller{
		strg: strg,
		log:  log,
	}
}

// Route generates routing for controller.
func (c *Controller) Route() *chi.Mux {
	r := chi.NewRouter()
	graphHandler, err := c.createHandler()
	if err != nil {
		c.log.Info("["+packageName+":Controller:Route] failed to create route: %v", err)
		r.HandleFunc("/", c.internalErrorHandler)
		return r
	}

	r.Handle("/", graphHandler)
	return r
}

// internalErrorHandler plug if handler was not created properly.
func (c *Controller) internalErrorHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "graphql currently unavailable", http.StatusInternalServerError)
}

// createHandler creates graph's handler.
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

// createSchema creates graph's schema.
func (c *Controller) createSchema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    c.createQueries(),
		Mutation: c.createMutations(),
	})
}

// createQueries creates Object with Select queries.
func (c *Controller) createQueries() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"selectPersons": c.selectPersonsQuery(),
		},
	})
}

// createMutations creates Object with mutable queries.
func (c *Controller) createMutations() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"insertPerson": c.insertPersonMutation(),
			"updatePerson": c.updatePersonMutation(),
			"deletePerson": c.deletePersonMutation(),
		},
	})
}

// selectPersonsQuery filter persons.
// filters: id(int), name(string), surname(string), patronymic(string),
// age(int), gender(string), country(string), page_num(int), page_size(int).
func (c *Controller) selectPersonsQuery() *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(personType),
		Args: getFieldConfigArgument([]string{
			argId,
			argName,
			argSurname,
			argPatronymic,
			argAge,
			argGender,
			argCountry,
			argPageNum,
			argPageSize,
		}),
		Resolve: c.selectPersonsResolver,
	}
}

// selectPersonsResolver resolver for selectPersonsQuery.
func (c *Controller) selectPersonsResolver(p graphql.ResolveParams) (interface{}, error) {
	foundPersons, err := c.strg.SelectPersons(p.Context, p.Args)
	if err != nil {
		c.log.Info("["+packageName+":Controller:selectPersonsResolver] find person failed: %v", err)
		return nil, fmt.Errorf("find persons failed: %w", err)
	}

	c.log.Info("["+packageName+":Controller:selectPersonsResolver] some persons were found by request filtering. count: '%d'", len(foundPersons))
	return foundPersons, nil
}

// insertPersonMutation filter persons.
// args:
//   - mandatory: name(string), surname(string), age(int), gender(string), country(string);
//   - optional: patronymic(string).
func (c *Controller) insertPersonMutation() *graphql.Field {
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
		Resolve: c.insertPersonResolver,
	}
}

// insertPersonResolver resolver for insertPersonMutation.
func (c *Controller) insertPersonResolver(p graphql.ResolveParams) (interface{}, error) {
	name, _ := p.Args[argName].(string)
	surname, _ := p.Args[argSurname].(string)
	patronymic, _ := p.Args[argPatronymic].(string)
	age, _ := p.Args[argAge].(int)
	gender, _ := p.Args[argGender].(string)
	country, _ := p.Args[argCountry].(string)

	newPerson := datastructs.PersonData{
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Age:        int64(age),
		Gender:     gender,
		Country:    country,
	}

	newPersonId, err := c.strg.AddPerson(p.Context, &newPerson)
	if err != nil {
		c.log.Info("["+packageName+":Controller:insertPersonResolver] create person failed: %w", err)
		return nil, fmt.Errorf("create person failed: %w", err)
	}

	newPerson.Id = newPersonId
	c.log.Info("["+packageName+":Controller:insertPersonResolver] person with id '%d' successfully created", newPersonId)
	return newPerson, nil
}

// updatePersonMutation modifies person partially by id.
// args:
//   - mandatory: id(int);
//   - optional: name(string), surname(string), patronymic(string), age(int), gender(string), country(string).
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

// updatePersonResolver resolver for updatePersonMutation.
func (c *Controller) updatePersonResolver(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args[argId].(int)
	updatedPerson, err := c.strg.UpdatePersonById(p.Context, int64(id), p.Args)
	if err != nil {
		c.log.Info("["+packageName+":Controller:updatePersonResolver] update person failed: %v", err)
		return nil, fmt.Errorf("update person failed: %w", err)
	}

	c.log.Info("["+packageName+":Controller:updatePersonResolver] person with id '%d' successfully updated", id)
	return updatedPerson, nil
}

// deletePersonMutation deletes person by id.
// args:
//   - mandatory: id(int).
func (c *Controller) deletePersonMutation() *graphql.Field {
	return &graphql.Field{
		Type: personType,
		Args: getFieldConfigArgument([]string{
			argId,
		}),
		Resolve: c.deletePersonResolver,
	}
}

// deletePersonResolver resolver for deletePersonMutation.
func (c *Controller) deletePersonResolver(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args[argId].(int)
	deletedPerson, err := c.strg.DeletePersonById(p.Context, int64(id))
	if err != nil {
		c.log.Info("["+packageName+":Controller:deletePersonResolver] delete person failed: %v", err)
		return nil, fmt.Errorf("delete person failed: %w", err)
	}

	c.log.Info("["+packageName+":Controller:deletePersonResolver] person with id '%d' successfully deleted", id)
	return deletedPerson, nil
}
