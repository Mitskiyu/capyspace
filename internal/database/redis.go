package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	rdb *redis.Client
}

func NewCache(rdb *redis.Client) *cache {
	return &cache{
		rdb: rdb,
	}
}

func ConnectRedis(password, host, port, db string) (*redis.Client, error) {
	connStr := fmt.Sprintf("redis://:%s@%s:%s/%s", password, host, port, db)
	opt, err := redis.ParseURL(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis url: %w", err)
	}

	rdb := redis.NewClient(opt)
	return rdb, nil
}

func PingRedis(ctx context.Context, rdb *redis.Client) error {
	if status := rdb.Ping(ctx); status.Err() != nil {
		return fmt.Errorf("failed to ping redis: %w", status.Err())
	}

	return nil
}
