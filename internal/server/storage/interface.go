package storage

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

type BaseStorage interface {
	AddPerson(ctx context.Context, data *datastructs.PersonData) (int64, error)
	GetPersons(ctx context.Context, filters map[string]interface{}, pageNum int64, pageSize int64) ([]datastructs.PersonData, error)
	DeletePersonById(ctx context.Context, personId int64) (*datastructs.PersonData, error)
	UpdatePersonById(ctx context.Context, personId int64, data *datastructs.PersonData) (int64, error)
	UpdatePersonByIdPartially(ctx context.Context, personId int64, values map[string]interface{}) (*datastructs.PersonData, error)
}
