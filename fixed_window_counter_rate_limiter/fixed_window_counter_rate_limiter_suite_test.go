package fixed_window_counter_ratelimiter_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFixedWindowCounterRateLimiter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FixedWindowCounterRateLimiter Suite")
}
