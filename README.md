# Learning Golang by Implementing Rate Limiter Algorithms

This project demonstrates various rate limiting algorithms implemented in Go, using Redis for distributed state management. Rate limiting is a crucial technique for controlling the amount of incoming and outgoing traffic to or from a network, application, or service.

## Overview

Rate limiting helps protect systems from abuse, ensures fair resource allocation, and maintains service availability. This repository contains implementations of various rate limiting algorithms, each with its own characteristics and use cases.

## Implemented Algorithms

### Token Bucket

The Token Bucket algorithm is a flexible rate limiting mechanism that works as follows:

- A bucket holds tokens, with a maximum capacity (bucket size)
- Tokens are added to the bucket at a constant rate (refill rate)
- When a request arrives, it takes one token from the bucket
- If the bucket has tokens, the request is allowed; otherwise, it's denied
- Allows for bursts of traffic up to the bucket capacity

**Key Properties:**

- Supports burst traffic handling
- Smooths out traffic over time
- Simple to understand and implement
- Configurable through bucket capacity and refill rate

**Implementation Details:**

- Uses Redis to store token counts and last refill timestamps
- Thread-safe and suitable for distributed environments
- Efficiently handles concurrent requests
- Implements a time-based token refill mechanism
- Automatically initializes new clients with a full bucket
- Prevents bucket overflow during token refill

**Testing:**

- Comprehensive test suite using Ginkgo and Gomega
- Tests cover all key scenarios:
  - Initial bucket filling for new clients
  - Token consumption for requests
  - Time-based token refilling
  - Bucket capacity enforcement
  - Empty bucket handling
  - Error handling for Redis failures

### Fixed Window Counter

The Fixed Window Counter algorithm is a straightforward rate limiting approach that works as follows:

- Time is divided into fixed windows (e.g., 1 second, 1 minute, etc.)
- Each client has a counter that tracks the number of requests in the current window
- When a new request arrives, the counter is incremented
- If the counter is below or at the limit, the request is allowed; otherwise, it's denied
- When a new time window starts, the counter resets to zero

**Key Properties:**

- Simple to understand and implement
- Low memory footprint (only stores one counter per client)
- Efficient for high-throughput systems
- Clear time boundaries for rate limits

**Implementation Details:**

- Uses Redis to store request counters with automatic expiration
- Counter keys automatically expire after the window duration
- Thread-safe and suitable for distributed environments
- Efficiently handles concurrent requests
- Automatically resets counters when a time window expires

**Testing:**

- Comprehensive test suite using Ginkgo and Gomega
- Tests cover all key scenarios:
  - New client handling
  - Request counting within limits
  - Limit enforcement
  - Window expiration and counter reset
  - Error handling for Redis failures

### Sliding Window Log (Coming Soon)

The Sliding Window Log algorithm will be implemented soon.

### Sliding Window Counter (Coming Soon)

The Sliding Window Counter algorithm will be implemented soon.

### Leaky Bucket (Coming Soon)

The Leaky Bucket algorithm will be implemented soon.

## Usage

### Prerequisites

- Go 1.25+
- Redis server

### Installation

```bash
# Clone the repository
git clone https://github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms.git
cd Learning-Golang-by-Implementing-Rate-Limiter-Algorithms

# Install dependencies
go mod download
go mod vendor
go mod tidy
```

### Example Usage

```go
package main

import (
    "context"
    "fmt"

    "github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"
    "github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/token_bucket_rate_limiter"
    "github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/fixed_window_counter_rate_limiter"
    "github.com/redis/go-redis/v9"
)

func main() {
    // Connect to Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Create Redis client wrapper
    rlRedisClient := rate_limiter.NewRedisClient(redisClient)

    // Example 1: Token Bucket Rate Limiter
    // Parameters: redisClient, bucketCapacity, refillRate
    tokenBucketRL := token_bucket_ratelimiter.NewTokenBucketRateLimiter(rlRedisClient, 10, 1)

    // Check if request is allowed
    clientID := "user123"
    if tokenBucketRL.LimitRequests(clientID) {
        fmt.Println("Request allowed")
    } else {
        fmt.Println("Request denied - rate limit exceeded")
    }

    // Example 2: Fixed Window Counter Rate Limiter
    // Parameters: redisClient, windowSize (in seconds), limit
    fixedWindowRL := fixed_window_counter_ratelimiter.NewFixedWindowCounterRateLimiter(rlRedisClient, 60, 5)

    // Check if request is allowed (5 requests per minute limit)
    if fixedWindowRL.LimitRequests(clientID) {
        fmt.Println("Request allowed")
    } else {
        fmt.Println("Request denied - rate limit exceeded")
    }
}
```

## Project Structure

```text
├── rate_limiter/
│   ├── mocks/
│   │   └── redis_client_mock.go          # Mock Redis client for testing
│   ├── rate_limiter.go           # Rate limiter interface definition
│   ├── redis_client.go           # Redis client wrapper implementation
│   └── redis_client_interface.go # Redis client interface definition
├── token_bucket_rate_limiter/
│   ├── token_bucket_rate_limiter.go      # Token Bucket implementation
│   ├── token_bucket_rate_limiter_suite_test.go  # Test suite setup
│   └── token_bucket_ratelimiter_test.go  # Comprehensive test cases
├── fixed_window_counter_rate_limiter/
│   ├── fixed_window_counter_rate_limiter.go      # Fixed Window Counter implementation
│   ├── fixed_window_counter_rate_limiter_suite_test.go  # Test suite setup
│   └── fixed_window_counter_ratelimiter_test.go  # Comprehensive test cases
└── [future_algorithm]/
    └── [future_algorithm].go  # Future algorithm implementations
```

## Contributing

Contributions are welcome! Feel free to submit a Pull Request.
