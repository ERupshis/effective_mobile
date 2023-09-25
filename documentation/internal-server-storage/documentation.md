package storage // import "github.com/erupshis/effective_mobile/internal/server/storage"

Package storage provides main storage interface and implementation.

TYPES

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
    BaseStorage main storage interface.

func Create(manager managers.BaseStorageManager, log logger.BaseLogger) BaseStorage
    Create base create method.

