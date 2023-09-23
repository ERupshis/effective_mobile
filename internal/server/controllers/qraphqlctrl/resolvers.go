package qraphqlctrl

import (
	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/graphql-go/graphql"
)

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

func (c *Controller) createPersonResolver(p graphql.ResolveParams) (interface{}, error) {
	// Create a new person
	name, _ := p.Args[argName].(string)
	surname, _ := p.Args[argSurname].(string)
	patronymic, _ := p.Args[argPatronymic].(string)
	age, _ := p.Args[argAge].(int)
	gender, _ := p.Args[argGender].(string)
	country, _ := p.Args[argCountry].(string)

	newPerson := datastructs.PersonData{
		Id:         int64(len(c.persons) + 1), //TODO: remove.
		Name:       name,
		Surname:    surname,
		Patronymic: patronymic,
		Age:        int64(age),
		Gender:     gender,
		Country:    country,
	}

	c.persons = append(c.persons, newPerson)
	return newPerson, nil
}

func (c *Controller) updatePersonResolver(p graphql.ResolveParams) (interface{}, error) {
	// Update a person by ID
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

func (c *Controller) deletePersonResolver(p graphql.ResolveParams) (interface{}, error) {
	// Delete a person by ID
	id, _ := p.Args[argId].(int64)
	for i, person := range c.persons {
		if person.Id == id {
			// Remove the person from the slice
			deletedPerson := c.persons[i]
			c.persons = append(c.persons[:i], c.persons[i+1:]...)
			return deletedPerson, nil
		}
	}
	return nil, nil
}
