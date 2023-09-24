// Package managers provides interface for storage manager.
package managers

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

// BaseStorageManager interface of storage manager that performs main work.
type BaseStorageManager interface {
	// AddPerson adds person in storage.
	AddPerson(ctx context.Context, data *datastructs.PersonData) (int64, error)
	// SelectPersons returns persons from storage satisfying filters conditions.
	SelectPersons(ctx context.Context, filters map[string]interface{}, pageNum int64, pageSize int64) ([]datastructs.PersonData, error)
	// DeletePersonById deletes person from storage by id.
	DeletePersonById(ctx context.Context, id int64) (int64, error)
	// UpdatePersonById updates person by id.
	UpdatePersonById(ctx context.Context, personId int64, values map[string]interface{}) (int64, error)
	// CheckConnection checks connection to storage.
	CheckConnection(ctx context.Context) (bool, error)
	// Close connection to storage. Should be called on close application stage.
	Close() error
}
