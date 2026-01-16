package redis

import (
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	RedisDB *redis.Client
}

func (rs *RedisStorage) InitRedisBD() {
	rs.RedisDB = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

func (rs *RedisStorage) GetUserPreferences(userID string, preferences map[string]interface{}) error {
	return nil
}

func (rs *RedisStorage) SetUserPreferences(userID string, preferences map[string]interface{}) error {
	return nil
}
