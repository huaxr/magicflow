// Author: huaxr
// Time:   2021/8/4 下午2:27
// Git:    huaxr

package accutil

import (
	"fmt"

	"github.com/huaxr/magicflow/component/logx"

	"sync"
	"sync/atomic"
)

//go:nosplit
func WrapF(func())

//go:nosplit
func wrapFN(string, func())

// Thread wrapper a goroutine and monitor its if exit.
// when using `go Thread("tmp-name", func(){f()})` please guarantee the
// the f is not closure inside.
// meanwhile `go Thread("tmp-name", f)` can use closure.
func Thread(name string, f func(), stop ...chan struct{}) {
	//defer func() {
	//	if err := recover(); err != nil {
	//		logx.L().Errorf("Thread.panic %v", err)
	//		// panic will not return to the previous PC address.
	//		// call exit here to monitor.
	//		exitMonitor(name)
	//	}
	//}()

	lock.Lock()
	if _, ok := monitor[name]; !ok {
		monitor[name] = &thread{
			goId:  getgid(),
			count: 1,
		}
	} else {
		atomic.AddInt32(&monitor[name].count, 1)
	}
	lock.Unlock()

	//log.Println("monitor.info", fmt.Sprintf("current thread: %+v", monitor[name]))
	wrapFN(name, f)
	return
}

func enter() {
	logx.L().Infof("WrapF.enter", "start goroutine: %d", getgid())
	return
}

func exit() {
	logx.L().Infof("WrapF.exit", "exit goroutine: %d", getgid())
	return
}

type thread struct {
	goId  int
	count int32
}

var (
	lock    = new(sync.Mutex)
	monitor = make(map[string]*thread)
)

func exitMonitor(name string) {
	lock.Lock()
	defer lock.Unlock()
	//logx.L().Infof("WrapF.exit", "exit goroutine: %d, name %v, current:%+v", getgid(), name, monitor)
	logx.L().Warnf("monitor exit %v", fmt.Sprintf("%s exit", name))

	// warning here!
	r := monitor[name]
	if r.count == 0 {
		delete(monitor, name)
	} else {
		atomic.AddInt32(&monitor[name].count, -1)
	}
	return
}
