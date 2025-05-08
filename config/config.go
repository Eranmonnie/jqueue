package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis(port string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     port,
		Password: "", // No password by default
		DB:       0,  // Default DB
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
}
