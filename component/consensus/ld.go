// Author: huaxr
// Time:   2021/6/12 下午12:36
// Git:    huaxr

package consensus

import "go.uber.org/atomic"

const leader = "__hitler"

// only leader can process
var leaderFlag = atomic.NewBool(false)

func MasterFunc(f func()) {
	if leaderFlag.Load() {
		f()
		return
	}
}
