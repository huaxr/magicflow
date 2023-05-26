// Author: huaxr
// Time:   2021/8/20 上午11:57
// Git:    huaxr

package registry

import (
	ctx2 "context"
	"fmt"
	"runtime"
	"time"

	"github.com/huaxr/magicflow/component/ticker"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/toolutil"
)

var HB *core.HeartBeat

func LoadHB() {
	HB = new(core.HeartBeat)
	inits := time.Now()

	f := func() {
		HB.CpuUsage = toolutil.Cpu()
		ms := toolutil.Mem()
		HB.MemIdle = fmt.Sprintf("%vM", int(ms.HeapIdle/(1024*1024)))
		HB.MemAlloc = fmt.Sprintf("%vM", int(ms.Alloc/(1024*1024)))
		HB.Alive = fmt.Sprintf("%.3fH", time.Now().Sub(inits).Hours())
		HB.OS = runtime.GOOS
		HB.Host = toolutil.GetIp()
	}
	tick := time.NewTicker(2 * time.Second)
	job := ticker.NewJob(ctx2.Background(), "local_state_hb", tick, f)
	ticker.GetManager().Register(job)
}
