package token_bucket_ratelimiter

import "github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"

var _ rate_limiter.RateLimiterInterface = (*TokenBucketRateLimiter)(nil)
