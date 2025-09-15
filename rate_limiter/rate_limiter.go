package rate_limiter

type RateLimiterInterface interface {
	LimitRequests(clientId string) bool
}