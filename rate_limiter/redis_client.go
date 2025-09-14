package rate_limiter

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx context.Context
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		client: client,
		ctx: context.Background(),
	}
}

func (r *RedisClient) GetCountAndLastRefill(keyCount, keyLastRefill string) (int64, int, error) {
	lastRefillStr, err := r.client.Get(r.ctx, keyLastRefill).Result()
	if err != nil && err != redis.Nil {
		return 0, 0, err
	}

	tokenCountStr, err := r.client.Get(r.ctx, keyCount).Result()
	if err != nil && err != redis.Nil {
		return 0, 0, err
	}

	var lastRefill int64
	var tokenCount int
	if lastRefillStr != "" {
		lastRefill, _ = strconv.ParseInt(lastRefillStr, 10, 64)
	}
	if tokenCountStr != "" {
		tokenCount, _ = strconv.Atoi(tokenCountStr)
	}

	return lastRefill, tokenCount, nil
}

func (r *RedisClient) SetCountAndLastRefill(keyCount, keyLastRefill string, tokenCount int, currentTime int64) error {
	if err := r.client.Set(r.ctx, keyLastRefill, strconv.FormatInt(currentTime, 10), 0).Err(); err != nil {
		return err
	}
	if err := r.client.Set(r.ctx, keyCount, strconv.Itoa(tokenCount), 0).Err(); err != nil {
		return err
	}
	return nil
}
