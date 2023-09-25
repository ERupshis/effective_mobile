package storage

import (
	"context"
	"fmt"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/helpers/requestshelper"
	"github.com/erupshis/effective_mobile/internal/server/storage/managers"
)

// storage main storage implementation. Consists of manager and logger.
type storage struct {
	manager managers.BaseStorageManager

	log logger.BaseLogger
}

// Create base create method.
func Create(manager managers.BaseStorageManager, log logger.BaseLogger) BaseStorage {
	return &storage{
		manager: manager,
		log:     log,
	}
}

// AddPerson adds person in storage. Perform person data validation and reject request if wrong. Returns new person id, otherwise -1.
func (s *storage) AddPerson(ctx context.Context, newPerson *datastructs.PersonData) (int64, error) {
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

// SelectPersons makes data selection from storage by filters, described in 'values'. Performs filtering 'values' and
// leave only valid filters for select. Returns array of persons. Also allows to paginate result.
func (s *storage) SelectPersons(ctx context.Context, values map[string]interface{}) ([]datastructs.PersonData, error) {
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

// DeletePersonById removes person from storage by id. Returns removed person data.
// Before person deleting search then in storage.
func (s *storage) DeletePersonById(ctx context.Context, id int64) (*datastructs.PersonData, error) {
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

// UpdatePersonById updates person selected by id. Applies filtering for input filter values. Returns person data.
func (s *storage) UpdatePersonById(ctx context.Context, id int64, values map[string]interface{}) (*datastructs.PersonData, error) {
	filteredValues := requestshelper.FilterValues(values)
	if len(filteredValues) == 0 {
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
