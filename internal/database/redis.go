package database

import (
	"context"
	"fmt"
	"time"

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
	if cmd := rdb.Ping(ctx); cmd.Err() != nil {
		return fmt.Errorf("failed to ping redis: %w", cmd.Err())
	}

	return nil
}

func (c *cache) SetSession(ctx context.Context, sessionId, userId string, exp time.Duration) error {
	return c.rdb.Set(ctx, sessionId, userId, exp).Err()
}

func (c *cache) GetSession(ctx context.Context, sessionId string) (string, error) {
	cmd := c.rdb.Get(ctx, sessionId)
	if cmd.Err() == redis.Nil {
		return "", fmt.Errorf("session does not exist")
	} else if cmd.Err() != nil {
		return "", cmd.Err()
	}

	return cmd.Val(), nil
}
