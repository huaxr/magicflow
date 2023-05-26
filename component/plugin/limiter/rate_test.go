// Author: huaxr
// Time:   2021/8/26 下午6:49
// Git:    huaxr

package limiter

import (
	"testing"
)

func TestNewRateLimiter(t *testing.T) {
	a := NewRateLimiter("a", 10)

	for {
		ok := a.Request()
		if ok {
			t.Log("success")
		} else {
			//t.Log("limited")
		}
	}

}
