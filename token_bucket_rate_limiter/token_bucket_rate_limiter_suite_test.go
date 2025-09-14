package token_bucket_ratelimiter_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTokenBucketRateLimiter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TokenBucketRateLimiter Suite")
}
