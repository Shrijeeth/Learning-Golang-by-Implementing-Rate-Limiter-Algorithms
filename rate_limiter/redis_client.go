package rate_limiter

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx context.Context
}

type ExpiryMode string
const (
	EXPIRY_MODE_DEFAULT ExpiryMode = ""
	EXPIRY_MODE_NX ExpiryMode = "NX"
	EXPIRY_MODE_XX ExpiryMode = "XX"
	EXPIRY_MODE_GT ExpiryMode = "GT"
	EXPIRY_MODE_LT ExpiryMode = "LT"
)

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

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Set(key string, value string) error {
	return r.client.Set(r.ctx, key, value, 0).Err()
}

func (r *RedisClient) Incr(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

func (r *RedisClient) Decr(key string) (int64, error) {
	return r.client.Decr(r.ctx, key).Result()
}

func (r *RedisClient) Expire(key string, duration time.Duration, expiryMode ExpiryMode) error {
	switch expiryMode {
		case EXPIRY_MODE_DEFAULT:
			return r.client.Expire(r.ctx, key, duration).Err()
		case EXPIRY_MODE_NX:
			return r.client.ExpireNX(r.ctx, key, duration).Err()
		case EXPIRY_MODE_XX:
			return r.client.ExpireXX(r.ctx, key, duration).Err()
		case EXPIRY_MODE_GT:
			return r.client.ExpireGT(r.ctx, key, duration).Err()
		case EXPIRY_MODE_LT:
			return r.client.ExpireLT(r.ctx, key, duration).Err()
		default:
			return errors.New("INVALID EXPIRY MODE")
	}
}

func (r *RedisClient) IncrWithExpiry(key string, duration time.Duration, expiryMode ExpiryMode) (int64, error) {
	var incrResult int64

	_, err := r.client.TxPipelined(r.ctx, func(pipe redis.Pipeliner) error {
		incrCmd := pipe.Incr(r.ctx, key)

		switch expiryMode {
			case EXPIRY_MODE_NX:
				pipe.ExpireNX(r.ctx, key, duration)
			case EXPIRY_MODE_XX:
				pipe.ExpireXX(r.ctx, key, duration)
			case EXPIRY_MODE_GT:
				pipe.ExpireGT(r.ctx, key, duration)
			case EXPIRY_MODE_LT:
				pipe.ExpireLT(r.ctx, key, duration)
			case EXPIRY_MODE_DEFAULT:
				pipe.Expire(r.ctx, key, duration)
			default:
				return errors.New("INVALID EXPIRY MODE")
		}
		incrResult = incrCmd.Val()

		return nil
	})

	return incrResult, err
}

