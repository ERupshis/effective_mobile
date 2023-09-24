// Package cache implements cache for storage.BaseStorage.
package cache

import (
	"context"
	"encoding/json"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/storage"
	"github.com/erupshis/effective_mobile/internal/server/storage/cache/manager"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// cache decoration class for storage.BaseStorage.
type cache struct {
	manager manager.BaseCacheManager
	storage storage.BaseStorage

	logger logger.BaseLogger
}

// Create creates cache-cover for storage.BaseStorage.
func Create(cacheManager manager.BaseCacheManager, baseStorage storage.BaseStorage, logger logger.BaseLogger) storage.BaseStorage {
	return &cache{
		manager: cacheManager,
		storage: baseStorage,
		logger:  logger,
	}
}

// AddPerson adds person in storage.
func (c *cache) AddPerson(ctx context.Context, data *datastructs.PersonData) (int64, error) {
	res, err := c.storage.AddPerson(ctx, data)
	if err == nil {
		if err = c.manager.Flush(ctx); err != nil {
			c.logger.Info("[cache:AddPerson] failed to flush cache %v", err)
		}
	}

	return res, err
}

// SelectPersons seeks and returns persons satisfying to filter.
func (c *cache) SelectPersons(ctx context.Context, filters map[string]interface{}) ([]datastructs.PersonData, error) {
	var res []datastructs.PersonData

	cachedVal, err := c.manager.Get(ctx, filters)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			c.logger.Info("[cache:SelectPersons] failed to get cacheValue with key: %v, error: %v", filters, err)
		}
	} else {
		err = json.Unmarshal(cachedVal, &res)
		if err == nil {
			c.logger.Info("[cache:SelectPersons] restored value from cache")
			return res, nil
		}
	}

	res, err = c.storage.SelectPersons(ctx, filters)
	if err == nil {
		jsonRes, err := json.Marshal(res)
		if err != nil {
			c.logger.Info("[cache:SelectPersons] failed to generate json result representation for storing: %v, error: %v", filters, err)
		} else {
			err = c.manager.Add(ctx, filters, jsonRes)
			if err != nil {
				c.logger.Info("[cache:SelectPersons] failed to add value in cache, value: %v, error: %v", jsonRes, err)
			}
		}
	}

	return res, err
}

// DeletePersonById deletes person from storage by id.
func (c *cache) DeletePersonById(ctx context.Context, id int64) (*datastructs.PersonData, error) {
	res, err := c.storage.DeletePersonById(ctx, id)
	if err == nil {
		if err = c.manager.Flush(ctx); err != nil {
			c.logger.Info("[cache:DeletePersonById] failed to flush cache %v", err)
		}
	}

	return res, err
}

// UpdatePersonById updates person data in storage by id.
func (c *cache) UpdatePersonById(ctx context.Context, id int64, values map[string]interface{}) (*datastructs.PersonData, error) {
	res, err := c.storage.UpdatePersonById(ctx, id, values)
	if err == nil {
		if err = c.manager.Flush(ctx); err != nil {
			c.logger.Info("[cache:UpdatePersonById] failed to flush cache %v", err)
		}
	}

	return res, err
}
