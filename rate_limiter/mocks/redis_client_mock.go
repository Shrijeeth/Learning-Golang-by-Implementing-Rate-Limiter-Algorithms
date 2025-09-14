package mocks

import (
	"errors"
)

// MockRedisClient implements the RedisClientInterface for testing
type MockRedisClient struct {
	GetCountAndLastRefillFunc func(keyCount, keyLastRefill string) (int64, int, error)
	SetCountAndLastRefillFunc func(keyCount, keyLastRefill string, tokenCount int, currentTime int64) error
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
