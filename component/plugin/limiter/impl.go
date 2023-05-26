// Author: huaxr
// Time:   2021/8/26 下午6:36
// Git:    huaxr

package limiter

type Limiter interface {
	// get the quota size
	GetQuota() int32
	GetType() RateType
	// try get a token
	Request() bool
}

type RateType int

const (
	GoRate RateType = iota + 1
	LeakyBucket
	SlidingWindow
	RedisTokenBucket
)

func LimiterFactory(rType RateType, app string, eps int32) Limiter {
	switch rType {
	case GoRate:
		return NewRateLimiter(app, eps)
	case LeakyBucket, SlidingWindow:
		panic("not implement yet")
	case RedisTokenBucket:

	default:
		return NewRateLimiter(app, eps)
	}
	return nil
}
