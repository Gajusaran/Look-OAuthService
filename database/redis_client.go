package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
)

func init() {

	redis_connecation_string := os.Getenv("REDIS_CONNECTION_STRING")

	Rdb = redis.NewClient(&redis.Options{
		Addr: redis_connecation_string,
		DB:   0, // total 16db is there we set each db for each usages
	})

	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully!")
}
