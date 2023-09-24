package manager

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/dgryski/go-farm"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/server/helpers/requestshelper"
	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client

	log logger.BaseLogger
}

func CreateRedis(ctx context.Context, dsn string, log logger.BaseLogger) (BaseCacheManager, error) {
	log.Info("[CreateRedis] open redis with settings: '%v'", dsn)

	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse redis DSN: %w", err)
	}

	client := redis.NewClient(opts)

	redisImpl := &Redis{
		client: client,
		log:    log,
	}

	ok, err := redisImpl.CheckConnection(ctx)
	if err != nil {
		return nil, fmt.Errorf("create redis: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("create redis: %w", err)
	}

	log.Info("[CreateRedis] successful")
	return redisImpl, nil
}

func (r *Redis) CheckConnection(ctx context.Context) (bool, error) {
	key := map[string]interface{}{"connection_check": ""}
	err := r.Add(ctx, key, "")
	if err != nil {
		return false, fmt.Errorf("check connection to redis: %w", err)
	}

	exists, err := r.Has(ctx, key)
	if err != nil {
		fmt.Println("Error:", err)
		return false, fmt.Errorf("check control key redis: %w", err)
	}

	return exists, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Add(ctx context.Context, key map[string]interface{}, val interface{}) error {
	err := r.client.Set(ctx, getKeyHash(key), val, 0).Err()
	if err != nil {
		return fmt.Errorf("redis add record: %w", err)
	}

	return nil
}

func (r *Redis) Has(ctx context.Context, key map[string]interface{}) (bool, error) {
	_, err := r.client.Get(ctx, getKeyHash(key)).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("redis check record prescense: %w", err)
	}
	return true, nil
}

func (r *Redis) Get(ctx context.Context, key map[string]interface{}) ([]byte, error) {
	val, err := r.client.Get(ctx, getKeyHash(key)).Result()
	if err == redis.Nil {
		return []byte{}, fmt.Errorf("missing record in redis: %w", err)
	} else if err != nil {
		return []byte{}, fmt.Errorf("redis get record: %w", err)
	} else {
		return []byte(val), nil
	}
}

func (r *Redis) Flush(ctx context.Context) error {
	_, err := r.client.FlushDB(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis flush: %w", err)
	}

	r.log.Info("[Redis:Flush] successful")
	return nil
}

type KeyVal struct {
	key   string
	value interface{}
}

func getKeyHash(values map[string]interface{}) string {
	coef := 37
	valueNum := 0

	var res uint64
	for _, val := range sortKeys(values) {
		multiplier := uint64(math.Pow(float64(coef), float64(valueNum)))
		res += multiplier * farm.Hash64([]byte(val.key))
		res += multiplier * farm.Hash64([]byte(requestshelper.ConvertQueryValueIntoString(val.value)))
	}

	sortKeys(values)
	return fmt.Sprintf("%x", res)
}

func sortKeys(values map[string]interface{}) []KeyVal {
	var res []KeyVal
	for k, v := range values {
		res = append(res, KeyVal{key: k, value: v})
	}

	sort.Slice(res, func(i, j int) bool { return res[i].key < res[j].key })
	return res
}
