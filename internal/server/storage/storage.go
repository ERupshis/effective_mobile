package storage

import (
	"context"
	"fmt"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/helpers/requestshelper"
	"github.com/erupshis/effective_mobile/internal/server/storage/managers"
)

type Storage struct {
	manager managers.BaseStorageManager

	log logger.BaseLogger
}

func Create(manager managers.BaseStorageManager, log logger.BaseLogger) Storage {
	return Storage{
		manager: manager,
		log:     log,
	}
}

func (s *Storage) AddPerson(ctx context.Context, newPerson *datastructs.PersonData) (int64, error) {
	_, err := requestshelper.IsPersonDataValid(newPerson, true)
	if err != nil {
		return -1, fmt.Errorf("resolve create person:%w", err)
	}

	newPersonId, err := s.manager.AddPerson(ctx, newPerson)
	if err != nil {
		return -1, fmt.Errorf("resolve create person: %w", err)
	}

	return newPersonId, nil
}

func (s *Storage) GetPersons(ctx context.Context, filters map[string]interface{}, pageNum int64, pageSize int64) ([]datastructs.PersonData, error) {
	return nil, nil
}

func (s *Storage) DeletePersonById(ctx context.Context, id int64) (*datastructs.PersonData, error) {
	personToDelete, err := s.manager.GetPersons(ctx, map[string]interface{}{"id": id}, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("delete person failed: %w", err)
	}

	affectedCount, err := s.manager.DeletePersonById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("delete person failed: %w", err)
	}

	if affectedCount == 0 {
		return nil, fmt.Errorf("person with id '%d' was not found", id)
	}

	return &personToDelete[0], nil
}
