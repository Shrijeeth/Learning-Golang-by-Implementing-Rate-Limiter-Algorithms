package rate_limiter

type RateLimiter interface {
	LimitRequests(clientId string) bool
}