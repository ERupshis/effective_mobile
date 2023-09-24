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

type Cache struct {
	manager manager.BaseCacheManager
	storage storage.BaseStorage

	logger logger.BaseLogger
}

func Create(cacheManager manager.BaseCacheManager, baseStorage storage.BaseStorage, logger logger.BaseLogger) storage.BaseStorage {
	return &Cache{
		manager: cacheManager,
		storage: baseStorage,
		logger:  logger,
	}
}

func (c *Cache) AddPerson(ctx context.Context, data *datastructs.PersonData) (int64, error) {
	res, err := c.storage.AddPerson(ctx, data)
	if err == nil {
		if err = c.manager.Flush(ctx); err != nil {
			c.logger.Info("[Cache:AddPerson] failed to flush cache %v", err)
		}
	}

	return res, err
}

func (c *Cache) SelectPersons(ctx context.Context, filters map[string]interface{}) ([]datastructs.PersonData, error) {
	var res []datastructs.PersonData

	cachedVal, err := c.manager.Get(ctx, filters)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			c.logger.Info("[Cache:SelectPersons] failed to get cacheValue with key: %v, error: %v", filters, err)
		}
	} else {
		err = json.Unmarshal(cachedVal, &res)
		if err == nil {
			c.logger.Info("[Cache:SelectPersons] restored value from cache")
			return res, nil
		}
	}

	res, err = c.storage.SelectPersons(ctx, filters)
	if err == nil {
		jsonRes, err := json.Marshal(res)
		if err != nil {
			c.logger.Info("[Cache:SelectPersons] failed to generate json result representation for storing: %v, error: %v", filters, err)
		} else {
			err = c.manager.Add(ctx, filters, jsonRes)
			if err != nil {
				c.logger.Info("[Cache:SelectPersons] failed to add value in cache, value: %v, error: %v", jsonRes, err)
			}
		}
	}

	return res, err
}

func (c *Cache) DeletePersonById(ctx context.Context, id int64) (*datastructs.PersonData, error) {
	res, err := c.storage.DeletePersonById(ctx, id)
	if err == nil {
		if err = c.manager.Flush(ctx); err != nil {
			c.logger.Info("[Cache:DeletePersonById] failed to flush cache %v", err)
		}
	}

	return res, err
}

func (c *Cache) UpdatePersonById(ctx context.Context, id int64, values map[string]interface{}) (*datastructs.PersonData, error) {
	res, err := c.storage.UpdatePersonById(ctx, id, values)
	if err == nil {
		if err = c.manager.Flush(ctx); err != nil {
			c.logger.Info("[Cache:UpdatePersonById] failed to flush cache %v", err)
		}
	}

	return res, err
}
