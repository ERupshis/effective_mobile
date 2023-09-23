package qraphqlctrl

import (
	"fmt"
	"net/http"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/helpers/requestshelper"
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

func (c *Controller) readPersonsResolver(p graphql.ResolveParams) (interface{}, error) {
	// Retrieve persons based on optional filter arguments
	//TODO: extend.
	name, _ := p.Args[argName].(string)
	surname, _ := p.Args[argSurname].(string)
	age, _ := p.Args[argAge].(int64)
	gender, _ := p.Args[argGender].(string)

	var filteredPersons = []datastructs.PersonData{}
	for _, person := range c.persons {
		// Check if each argument matches the person's data
		if (name == "" || person.Name == name) &&
			(surname == "" || person.Surname == surname) &&
			(age == 0 || person.Age == age) &&
			(gender == "" || person.Gender == gender) {
			filteredPersons = append(filteredPersons, person)
		}
	}
	return filteredPersons, nil
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

func (c *Controller) createPersonResolver(p graphql.ResolveParams) (interface{}, error) {
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

	_, err := requestshelper.IsPersonDataValid(&newPerson, true)
	if err != nil {
		c.log.Info("["+packageName+":Controller:createPersonResolver] data validation failed: %v", err)
		return newPerson, fmt.Errorf("resolve create person:%w", err)
	}

	newPersonId, err := c.strg.AddPerson(p.Context, &newPerson)
	if err != nil {
		c.log.Info("["+packageName+":Controller:createPersonResolver] cannot process: %v", err)
		return newPerson, fmt.Errorf("resolve create person: %w", err)
	}

	c.log.Info("["+packageName+":Controller:createPersonResolver] person successfully added with id '%d'", newPersonId)
	newPerson.Id = newPersonId
	return newPerson, nil
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

func (c *Controller) updatePersonResolver(p graphql.ResolveParams) (interface{}, error) {
	//TODO: need to support partial update.
	id, _ := p.Args["argId"].(int64)
	for i, person := range c.persons {
		if person.Id == id {
			if name, ok := p.Args[argName].(string); ok {
				c.persons[i].Name = name
			}
			if surname, ok := p.Args[argSurname].(string); ok {
				c.persons[i].Surname = surname
			}
			if patronymic, ok := p.Args[argPatronymic].(string); ok {
				c.persons[i].Patronymic = patronymic
			}
			if age, ok := p.Args[argAge].(int64); ok {
				c.persons[i].Age = age
			}
			if gender, ok := p.Args[argGender].(string); ok {
				c.persons[i].Gender = gender
			}
			if country, ok := p.Args[argCountry].(string); ok {
				c.persons[i].Country = country
			}
			return c.persons[i], nil
		}
	}
	return nil, nil
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

func (c *Controller) deletePersonResolver(p graphql.ResolveParams) (interface{}, error) {
	// Delete a person by ID
	id, _ := p.Args[argId].(int)

	personToDelete, err := c.strg.GetPersons(p.Context, map[string]interface{}{"id": id}, 0, 0)
	if err != nil {
		c.log.Info("["+packageName+":Controller:deletePersonByIdHandler] person was not found in storage by id '%s': %v", id, err)
		return nil, fmt.Errorf("delete person failed: %w", err)
	}

	affectedCount, err := c.strg.DeletePersonById(p.Context, int64(id))
	if err != nil {
		c.log.Info("["+packageName+":Controller:deletePersonByIdHandler] person id is not valid: %v", err)
		return nil, fmt.Errorf("delete person failed: %w", err)
	}

	if affectedCount == 0 {
		c.log.Info("["+packageName+":Controller:deletePersonResolver] request has no effect with id '%d'", id)
		return nil, fmt.Errorf("person with id '%d' was not found", id)
	}

	c.log.Info("["+packageName+":Controller:deletePersonResolver] person with id '%d' successfully deleted", id)
	return personToDelete[0], nil
}
