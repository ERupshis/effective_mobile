package storage

import (
	"context"

	"github.com/erupshis/effective_mobile/internal/datastructs"
)

type BaseStorageManager interface {
	AddPersonData(ctx context.Context, data *datastructs.PersonData) error
	GetPersonsData(ctx context.Context) ([]datastructs.PersonData, error)
	DeletePersonDataById(ctx context.Context, personId int64) error
	UpdatePersonDataById(ctx context.Context, personId int64) error
	CheckConnection(ctx context.Context) (bool, error)
	Close() error
}
