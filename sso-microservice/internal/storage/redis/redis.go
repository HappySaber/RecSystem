package redisStorage

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	RedisDB *redis.Client
}

type RedisConfig struct {
	Addr        string
	Password    string
	User        string
	DB          int
	MaxRetries  int
	DialTimeout time.Duration
	Timeout     time.Duration
}

func NewConfig() (*RedisConfig, error) {
	const op = "storage.redis.NewConfig"

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("%s: REDIS_ADDR is required", op)
	}

	password := os.Getenv("REDIS_PASSWORD")
	user := os.Getenv("REDIS_USER")

	db, err := getEnvInt("REDIS_DB", 0)
	if err != nil {
		return nil, fmt.Errorf("%s: REDIS_DB: %w", op, err)
	}

	maxRetries, err := getEnvInt("REDIS_MAX_RETRIES", 3)
	if err != nil {
		return nil, fmt.Errorf("%s: REDIS_MAX_RETRIES: %w", op, err)
	}

	dialTimeout, err := getEnvDuration("REDIS_DIAL_TIMEOUT", "5s")
	if err != nil {
		return nil, fmt.Errorf("%s: REDIS_DIAL_TIMEOUT: %w", op, err)
	}

	timeout, err := getEnvDuration("REDIS_TIMEOUT", "3s")
	if err != nil {
		return nil, fmt.Errorf("%s: REDIS_TIMEOUT: %w", op, err)
	}

	return &RedisConfig{
		Addr:        addr,
		Password:    password,
		User:        user,
		DB:          db,
		MaxRetries:  maxRetries,
		DialTimeout: dialTimeout,
		Timeout:     timeout,
	}, nil
}

func getEnvInt(key string, def int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return n, nil
}

func getEnvDuration(key string, def string) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" {
		v = def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return d, nil
}

func NewClient(ctx context.Context) (*redis.Client, error) {
	const op = "storage.redis.NewClient"
	cfg, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := db.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

func (rs *RedisStorage) SaveToken(ctx context.Context, token, userID string, ttl time.Duration) error {
	const op = "storage.redis.SaveToken"

	if err := rs.RedisDB.Set(ctx, "session:"+token, userID, ttl).Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (rs *RedisStorage) Logout(ctx context.Context, token string) (bool, error) {
	const op = "storage.redis.Logout"

	if err := rs.RedisDB.Del(ctx, "session:"+token).Err(); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}
