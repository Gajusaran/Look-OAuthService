package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Rdb *redis.Client
)

func init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0, // total 16db is there we set each db for each usages
	})

	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully!")
}
