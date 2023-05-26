// Author: XinRui Hua
// Time:   2022/4/26 下午10:40
// Git:    huaxr

package metric

import (
	"sync"
	"sync/atomic"
)

// guage pie by apps
type appmetricx struct {
	apps    []int32
	indexes map[string]int
	lock    sync.RWMutex
}

var (
	app = &appmetricx{
		apps:    make([]int32, 0),
		indexes: make(map[string]int),
		lock:    sync.RWMutex{},
	}
)

func RegisterApp(appName string) {
	app.lock.Lock()
	defer app.lock.Unlock()

	index := len(app.apps)
	app.apps = append(app.apps, 0)
	app.indexes[appName] = index - 1
}

func MetricApp(appName string) {
	app.lock.RLock()
	defer app.lock.RUnlock()
	index := app.indexes[appName]
	atomic.AddInt32(&app.apps[index], 1)
}

func CleanApp(appName string) {
	app.lock.RLock()
	defer app.lock.RUnlock()
	index := app.indexes[appName]
	atomic.SwapInt32(&tags[index].count, 0)
}
