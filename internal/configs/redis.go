package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	rdbHost := os.Getenv("RDBHOST")
	rdbPort := os.Getenv("RDBPORT")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rdbHost, rdbPort),
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	} else {
		log.Println("RDB Connected")
	}

	return rdb, nil
}
