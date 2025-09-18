package mocks

import (
	"errors"
	"time"

	"github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"
)

// MockRedisClient implements the RedisClientInterface for testing
type MockRedisClient struct {
	GetCountAndLastRefillFunc func(keyCount, keyLastRefill string) (int64, int, error)
	SetCountAndLastRefillFunc func(keyCount, keyLastRefill string, tokenCount int, currentTime int64) error
	GetFunc                   func(key string) (string, error)
	SetFunc                   func(key string, value string) error
	IncrFunc                  func(key string) (int64, error)
	DecrFunc                  func(key string) (int64, error)
	ExpireFunc                func(key string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) error
	IncrWithExpiryFunc        func(key string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error)
	HGetAllFunc               func(key string) (map[string]string, error)
	HIncrByWithExpiryFunc     func(key string, value string, increment int64, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error)
	HLenFunc                  func(key string) (int64, error)
	HSetWithExpiryFunc        func(key string, value string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error)
}

// NewMockRedisClient creates a new mock Redis client
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{}
}

// GetCountAndLastRefill overrides the RedisClient method for testing
func (m *MockRedisClient) GetCountAndLastRefill(keyCount, keyLastRefill string) (int64, int, error) {
	if m.GetCountAndLastRefillFunc != nil {
		return m.GetCountAndLastRefillFunc(keyCount, keyLastRefill)
	}
	return 0, 0, errors.New("GetCountAndLastRefill not implemented")
}

// SetCountAndLastRefill overrides the RedisClient method for testing
func (m *MockRedisClient) SetCountAndLastRefill(keyCount, keyLastRefill string, tokenCount int, currentTime int64) error {
	if m.SetCountAndLastRefillFunc != nil {
		return m.SetCountAndLastRefillFunc(keyCount, keyLastRefill, tokenCount, currentTime)
	}
	return errors.New("SetCountAndLastRefill not implemented")
}

// Get overrides the RedisClient method for testing
func (m *MockRedisClient) Get(key string) (string, error) {
	if m.GetFunc != nil {
		return m.GetFunc(key)
	}
	return "", errors.New("Get not implemented")
}

// Set overrides the RedisClient method for testing
func (m *MockRedisClient) Set(key string, value string) error {
	if m.SetFunc != nil {
		return m.SetFunc(key, value)
	}
	return errors.New("Set not implemented")
}

// Incr overrides the RedisClient method for testing
func (m *MockRedisClient) Incr(key string) (int64, error) {
	if m.IncrFunc != nil {
		return m.IncrFunc(key)
	}
	return 0, errors.New("Incr not implemented")
}

// Decr overrides the RedisClient method for testing
func (m *MockRedisClient) Decr(key string) (int64, error) {
	if m.DecrFunc != nil {
		return m.DecrFunc(key)
	}
	return 0, errors.New("Decr not implemented")
}

// Expire overrides the RedisClient method for testing
func (m *MockRedisClient) Expire(key string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) error {
	if m.ExpireFunc != nil {
		return m.ExpireFunc(key, duration, expiryMode)
	}
	return errors.New("Expire not implemented")
}

func (m *MockRedisClient) IncrWithExpiry(key string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error) {
	if m.IncrWithExpiryFunc != nil {
		return m.IncrWithExpiryFunc(key, duration, expiryMode)
	}
	return 0, errors.New("IncrWithExpiry not implemented")
}

func (m *MockRedisClient) HGetAll(key string) (map[string]string, error) {
	if m.HGetAllFunc != nil {
		return m.HGetAllFunc(key)
	}
	return nil, errors.New("HGetAll not implemented")
}

func (m *MockRedisClient) HIncrByWithExpiry(key string, value string, increment int64, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error) {
	if m.HIncrByWithExpiryFunc != nil {
		return m.HIncrByWithExpiryFunc(key, value, increment, duration, expiryMode)
	}
	return 0, errors.New("HIncrByWithExpiry not implemented")
}

func (m *MockRedisClient) HLen(key string) (int64, error) {
	if m.HLenFunc != nil {
		return m.HLenFunc(key)
	}
	return 0, errors.New("HLen not implemented")
}

func (m *MockRedisClient) HSetWithExpiry(key string, value string, duration time.Duration, expiryMode rate_limiter.ExpiryMode) (int64, error) {
	if m.HSetWithExpiryFunc != nil {
		return m.HSetWithExpiryFunc(key, value, duration, expiryMode)
	}
	return 0, errors.New("HSetWithExpiry not implemented")
}
