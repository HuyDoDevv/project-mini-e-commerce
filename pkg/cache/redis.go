package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheService struct {
	ctx context.Context
	rdb *redis.Client
}

func NewRedisCacheService(rdb *redis.Client) *RedisCacheService {
	return &RedisCacheService{
		ctx: context.Background(),
		rdb: rdb,
	}
}

func (rcs *RedisCacheService) Get(key string, dest any) error {
	data, err := rcs.rdb.Get(rcs.ctx, key).Result()
	if err != nil {
		return err
	}
	if errors.Is(err, redis.Nil) {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func (rcs *RedisCacheService) Set(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rcs.rdb.Set(rcs.ctx, key, data, ttl).Err()
}

func (rcs *RedisCacheService) Clear(patten string) error {
	cursor := uint64(0)
	for {
		keys, nextCusor, err := rcs.rdb.Scan(rcs.ctx, cursor, patten, 2).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			rcs.rdb.Del(rcs.ctx, keys...)
		}
		cursor = nextCusor

		if cursor == 0 {
			break
		}
	}
	return nil
}
