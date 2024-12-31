package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis adresse par défaut
		Password: "",               // pas de mot de passe par défaut
		DB:       0,                // base de données par défaut
	})

	// Test de la connexion
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("could not connect to redis: %v", err)
	}

	fmt.Println("🚀 Redis connected successfully")
	return nil
}
