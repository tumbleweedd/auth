package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/tumbleweedd/svc/auth_service/internal/config"
	"time"
)

func NewRedisClient(cfg *config.RedisConfig) *redis.Client {
	opts := &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
		Password:     cfg.Password,
		DB:           cfg.DB,
	}

	return redis.NewClient(opts)
}
