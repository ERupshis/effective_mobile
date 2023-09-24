// Package storage provides main storage interface and implementation.
package storage

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

// BaseStorage main storage interface.
type BaseStorage interface {
	// AddPerson adds person in storage.
	AddPerson(ctx context.Context, data *datastructs.PersonData) (int64, error)
	// SelectPersons seeks and returns persons satisfying to filter.
	SelectPersons(ctx context.Context, filters map[string]interface{}) ([]datastructs.PersonData, error)
	// DeletePersonById deletes person from storage by id.
	DeletePersonById(ctx context.Context, id int64) (*datastructs.PersonData, error)
	// UpdatePersonById updates person data in storage by id.
	UpdatePersonById(ctx context.Context, id int64, values map[string]interface{}) (*datastructs.PersonData, error)
}
