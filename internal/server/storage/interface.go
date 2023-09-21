package storage

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

type BaseStorageManager interface {
	AddPersonData(ctx context.Context, data *datastructs.PersonData) error
	GetPersonsData(ctx context.Context, filters map[string]string, pageLimit int64, offset int64) ([]datastructs.PersonData, error)
	DeletePersonDataById(ctx context.Context, personId int64) error
	UpdatePersonDataById(ctx context.Context, personId int64, data *datastructs.PersonData) error
	UpdatePersonDataByIdPartially(ctx context.Context, personId int64, values map[string]string) error
	CheckConnection(ctx context.Context) (bool, error)
	Close() error
}
