package qraphqlctrl

import (
	"fmt"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/storage"
	"github.com/go-chi/chi/v5"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

const packageName = "graphqlctrl"

var personType = graphql.NewObject(graphql.ObjectConfig{
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

type Controller struct {
	persons []datastructs.PersonData

	strg storage.BaseStorageManager
	log  logger.BaseLogger
}

func Create(strg storage.BaseStorageManager, log logger.BaseLogger) *Controller {
	ctrl := &Controller{
		strg: strg,
		log:  log,
	}
	ctrl.persons = []datastructs.PersonData{
		{Id: 1, Name: "John", Surname: "Doe", Patronymic: "Smith", Age: 30, Gender: "Male", Country: "USA"},
		{Id: 2, Name: "Jane", Surname: "Doe", Patronymic: "Johnson", Age: 25, Gender: "Female", Country: "USA"},
	}

	return &Controller{
		strg: strg,
		log:  log,
	}
}

//func Create(strg storage.BaseStorageManager, log logger.BaseLogger) *Controller {
//	return &Controller{
//		strg: strg,
//		log: log,
//	}
//}

func (c *Controller) Route() *chi.Mux {
	r := chi.NewRouter()
	graphHandler, err := c.createHandler()
	if err != nil {
		c.log.Info("["+packageName+":Controller:Route] failed to create route: %v", err)
		return r
	}

	r.Handle("/", graphHandler)
	return r
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
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"surname": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"age": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"gender": &graphql.ArgumentConfig{
				Type: graphql.String,
			},

			//TODO: extend to all params as filters
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// Retrieve persons based on optional filter arguments
			//TODO: extend.
			name, _ := p.Args["name"].(string)
			surname, _ := p.Args["surname"].(string)
			age, _ := p.Args["age"].(int64)
			gender, _ := p.Args["gender"].(string)

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
		},
	}
}

func (c *Controller) createPersonMutation() *graphql.Field {
	return &graphql.Field{
		Type: personType,
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"surname": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"patronymic": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"age": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"gender": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"country": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// Create a new person
			name, _ := p.Args["name"].(string)
			surname, _ := p.Args["surname"].(string)
			patronymic, _ := p.Args["patronymic"].(string)
			age, _ := p.Args["age"].(int)
			gender, _ := p.Args["gender"].(string)
			country, _ := p.Args["country"].(string)

			newPerson := datastructs.PersonData{
				Id:         int64(len(c.persons) + 1),
				Name:       name,
				Surname:    surname,
				Patronymic: patronymic,
				Age:        int64(age),
				Gender:     gender,
				Country:    country,
			}

			c.persons = append(c.persons, newPerson)
			return newPerson, nil
		},
	}
}

func (c *Controller) updatePersonMutation() *graphql.Field {
	return &graphql.Field{
		Type: personType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"surname": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"patronymic": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"age": &graphql.ArgumentConfig{
				Type: graphql.Int,
			},
			"gender": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"country": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// Update a person by ID
			//TODO: need to support partial update.
			id, _ := p.Args["id"].(int64)
			for i, person := range c.persons {
				if person.Id == id {
					if name, ok := p.Args["name"].(string); ok {
						c.persons[i].Name = name
					}
					if surname, ok := p.Args["surname"].(string); ok {
						c.persons[i].Surname = surname
					}
					if patronymic, ok := p.Args["patronymic"].(string); ok {
						c.persons[i].Patronymic = patronymic
					}
					if age, ok := p.Args["age"].(int64); ok {
						c.persons[i].Age = age
					}
					if gender, ok := p.Args["gender"].(string); ok {
						c.persons[i].Gender = gender
					}
					if country, ok := p.Args["country"].(string); ok {
						c.persons[i].Country = country
					}
					return c.persons[i], nil
				}
			}
			return nil, nil
		},
	}
}

func (c *Controller) deletePersonMutation() *graphql.Field {
	return &graphql.Field{
		Type: personType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// Delete a person by ID
			id, _ := p.Args["id"].(int64)
			for i, person := range c.persons {
				if person.Id == id {
					// Remove the person from the slice
					deletedPerson := c.persons[i]
					c.persons = append(c.persons[:i], c.persons[i+1:]...)
					return deletedPerson, nil
				}
			}
			return nil, nil
		},
	}
}
