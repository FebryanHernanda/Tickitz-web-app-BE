package utils

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func GetCache(ctx context.Context, rdb *redis.Client, key string, data interface{}) error {
	// Check if redis server down
	if rdb == nil {
		log.Println("Redis server unavailable")
		return nil
	}

	cmd := rdb.Get(ctx, key)
	// check if cache miss return nil
	if cmd.Err() == redis.Nil {
		log.Println("Cache miss for key:", key)
		return nil
	} else if cmd.Err() != nil {
		log.Println("Redis error: ", cmd.Err())
		return cmd.Err()
	}

	// parse JSON data
	if err := json.Unmarshal([]byte(cmd.Val()), data); err != nil {
		log.Println("Redis unmarshal error: ", cmd.Err())
		return err
	}

	log.Println("Cache hit for key:", key)
	return nil
}

func SetCache(ctx context.Context, rdb *redis.Client, key string, data interface{}, ttl time.Duration) error {
	// Check if redis server down
	if rdb == nil {
		log.Println("Redis server unavailable")
		return nil
	}

	// Parse all type data to JSON
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Redis marshal error: ", err)
		return err
	}

	// Set data to redis
	if err := rdb.Set(ctx, key, bytes, ttl).Err(); err != nil {
		log.Println("Redis set error: ", err)
		return err
	}

	return nil
}

func InvalidateCache(ctx context.Context, rdb *redis.Client, prefixes []string) error {
	if rdb == nil {
		log.Println("Redis server unavailable")
		return nil
	}

	for _, prefix := range prefixes {
		iter := rdb.Scan(ctx, 0, prefix+"*", 0).Iterator()
		for iter.Next(ctx) {
			err := rdb.Del(ctx, iter.Val()).Err()
			if err != nil {
				log.Println("Redis delete error: ", err)
				return err
			}
		}

		if err := iter.Err(); err != nil {
			log.Println("Redis scan error : ", err)
			return err
		}
	}

	return nil
}
