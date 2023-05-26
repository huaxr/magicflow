// Author: huaxr
// Time:   2021/8/26 下午6:38
// Git:    huaxr

package limiter

import (
	"golang.org/x/time/rate"
	"time"
)

type RateLimiter struct {
	limit *rate.Limiter
	app   string
	quota int32
}

func (rl *RateLimiter) GetQuota() int32 {
	return rl.quota
}

func (rl *RateLimiter) GetType() RateType {
	return GoRate
}

func (rl *RateLimiter) Request() bool {
	return rl.limit.Allow()
}

func NewRateLimiter(app string, eps int32) Limiter {
	// eventy 1 millisecond put a token in this bucket
	limit := rate.Every(time.Duration(1e6/eps) * time.Microsecond)
	// the second param is the bucket size.
	limiter := rate.NewLimiter(limit, int(eps)*2)

	rl := new(RateLimiter)
	rl.limit = limiter
	rl.app = app
	rl.quota = eps
	return rl
}
