package sliding_window_counter_rate_limiter

import "github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"

var _ rate_limiter.RateLimiterInterface = (*SlidingWindowCounterRateLimiter)(nil)
