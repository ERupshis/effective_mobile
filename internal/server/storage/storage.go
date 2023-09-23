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

func Create(manager managers.BaseStorageManager, log logger.BaseLogger) BaseStorage {
	return &Storage{
		manager: manager,
		log:     log,
	}
}

func (s *Storage) AddPerson(ctx context.Context, newPerson *datastructs.PersonData) (int64, error) {
	_, err := requestshelper.IsPersonDataValid(newPerson, true)
	if err != nil {
		return -1, fmt.Errorf("storage: create: %w", err)
	}

	newPersonId, err := s.manager.AddPerson(ctx, newPerson)
	if err != nil {
		return -1, fmt.Errorf("storage: create: %w", err)
	}

	return newPersonId, nil
}

func (s *Storage) SelectPersons(ctx context.Context, values map[string]interface{}) ([]datastructs.PersonData, error) {
	filters := requestshelper.FilterValues(values)

	if _, ok := values["id"]; ok {
		filters["id"] = values["id"]
	}

	filters, pageNum, pageSize := requestshelper.FilterPageNumAndPageSize(filters)
	if pageSize < 0 {
		return nil, fmt.Errorf("storage: get: negative page size")
	}

	if pageNum < 0 {
		return nil, fmt.Errorf("storage: get: negative page num")
	}

	selectedPersons, err := s.manager.SelectPersons(ctx, filters, pageNum, pageSize)
	if err != nil {
		return nil, fmt.Errorf("storage: get: cannot process: %w", err)
	}

	if len(selectedPersons) == 0 {
		return nil, fmt.Errorf("storage: get: any person in storage does not meet the filtering conditions")
	}

	return selectedPersons, nil
}

func (s *Storage) DeletePersonById(ctx context.Context, id int64) (*datastructs.PersonData, error) {
	personToDelete, err := s.manager.SelectPersons(ctx, map[string]interface{}{"id": id}, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("storage: delete: %w", err)
	}

	affectedCount, err := s.manager.DeletePersonById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("storage: delete: %w", err)
	}

	if affectedCount == 0 {
		return nil, fmt.Errorf("storage: delete: person with id '%d' was not found", id)
	}

	return &personToDelete[0], nil
}

func (s *Storage) UpdatePersonById(ctx context.Context, id int64, values map[string]interface{}) (*datastructs.PersonData, error) {
	filteredValues := requestshelper.FilterValues(values)
	if len(values) == 0 {
		return nil, fmt.Errorf("storage: update: missing correct filters in request")
	}

	affectedCount, err := s.manager.UpdatePersonById(ctx, id, filteredValues)
	if err != nil {
		return nil, fmt.Errorf("storage: update: cannot process: %w", err)
	}

	if affectedCount == 0 {
		return nil, fmt.Errorf("storage: update: has no effect with id '%d'", id)
	}

	updatedPerson, err := s.manager.SelectPersons(ctx, map[string]interface{}{"id": id}, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("storage: update: failed to find updated person with id '%d'", id)
	}

	return &updatedPerson[0], nil
}
