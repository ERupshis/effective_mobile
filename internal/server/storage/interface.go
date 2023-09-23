package storage

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

type BaseStorage interface {
	AddPerson(ctx context.Context, data *datastructs.PersonData) (int64, error)
	SelectPersons(ctx context.Context, filters map[string]interface{}) ([]datastructs.PersonData, error)
	DeletePersonById(ctx context.Context, id int64) (*datastructs.PersonData, error)
	UpdatePersonById(ctx context.Context, id int64, values map[string]interface{}) (*datastructs.PersonData, error)
}
