// Author: XinRui Hua
// Time:   2022/4/22 下午2:31
// Git:    huaxr

package vars

import (
	"sync/atomic"

	"github.com/huaxr/magicflow/component/logx"
)

var (
	traceCount int

	syncCount int

	ackCount int

	residualTask int32
)

func GetTraceCount() int { return traceCount }

func SetTraceCount(i int) { traceCount = i }

func GetSyncCount() int { return syncCount }

func SetSyncCount(i int) { syncCount = i }

func GetAckCount() int { return ackCount }

func SetAckCount(i int) { ackCount = i }

func IncResidualTask() { atomic.AddInt32(&residualTask, 1) }

func DecResidualTask() {
	// consider concurrency conditions when a sync message timeout
	// but execute finished at the same time(on the multi branches),
	// so Complete and Exception will calling the dec twice.
	if atomic.LoadInt32(&residualTask) < 1 {
		logx.L().Errorf("residualTask could not less than zero")
		return
	}
	atomic.AddInt32(&residualTask, -1)
}

func GetResidualTask() int32 {
	return residualTask
}
