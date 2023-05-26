// Author: huaxr
// Time:   2021/12/22 下午5:26
// Git:    huaxr

package toolutil

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"runtime"
	"time"
)

func Cpu() string {
	percent, _ := cpu.Percent(time.Second, false)
	x := fmt.Sprintf("%.1f", percent[0])
	return fmt.Sprintf("%v%%", x)
}

func Mem() runtime.MemStats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return ms
}
