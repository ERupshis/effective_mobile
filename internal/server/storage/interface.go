package storage

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

type BaseStorageManager interface {
	AddPerson(ctx context.Context, data *datastructs.PersonData) error
	GetPersons(ctx context.Context, filters map[string]string, page int64, pageSize int64) ([]datastructs.PersonData, error)
	DeletePersonById(ctx context.Context, personId int64) (int64, error)
	UpdatePersonById(ctx context.Context, personId int64, data *datastructs.PersonData) (int64, error)
	UpdatePersonByIdPartially(ctx context.Context, personId int64, values map[string]string) error
	CheckConnection(ctx context.Context) (bool, error)
	Close() error
}
