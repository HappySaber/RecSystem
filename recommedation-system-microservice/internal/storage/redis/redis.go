package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

type RedisStorage struct {
	RedisDB *redis.Client
}

type RedisConfig struct {
	Addr        string
	Password    string
	User        string `yaml:"user"`
	DB          int
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
}

func NewConfig() (*RedisConfig, error) {
	const op = "storage.redis.NewConfig"

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "local.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("%s: read config: %w", op, err)
	}

	cfg := RedisConfig{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("%s: unmarshal yaml: %w", op, err)
	}

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("%s: REDIS_ADDR is required", op)
	}

	dbStr := os.Getenv("REDIS_DB")
	if dbStr == "" {
		dbStr = "0"
	}

	dbNum, err := strconv.Atoi(dbStr)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid REDIS_DB: %w", op, err)
	}

	cfg.Addr = addr
	cfg.Password = os.Getenv("REDIS_PASSWORD")
	cfg.DB = dbNum

	return &cfg, nil
}

func NewClient(ctx context.Context, cfg RedisConfig) (*redis.Client, error) {
	const op = "storage.redis.NewClient"

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

func (rs *RedisStorage) GetUserPreferences(userID string, preferences map[string]interface{}) error {
	return nil
}

func (rs *RedisStorage) SetUserPreferences(userID string, preferences map[string]interface{}) error {
	return nil
}

func (rs *RedisStorage) ResetUserPreferences(userID string) error {
	return nil
}
