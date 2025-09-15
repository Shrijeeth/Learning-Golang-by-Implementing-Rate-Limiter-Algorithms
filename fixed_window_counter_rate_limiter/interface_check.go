package fixed_window_counter_ratelimiter

import "github.com/Shrijeeth/Learning-Golang-by-Implementing-Rate-Limiter-Algorithms/rate_limiter"

var _ rate_limiter.RateLimiterInterface = (*FixedWindowCounterRateLimiter)(nil)
