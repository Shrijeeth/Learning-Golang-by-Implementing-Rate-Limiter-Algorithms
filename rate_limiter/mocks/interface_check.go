package mocks

import (
	"github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"
)

// This is just a compile-time check to ensure MockRedisClient implements RedisClientInterface
var _ rate_limiter.RedisClientInterface = (*MockRedisClient)(nil)
