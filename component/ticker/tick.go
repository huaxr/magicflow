// Author: huaxr
// Time: 2022/7/5 10:48 上午
// Git: huaxr

package ticker

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"github.com/huaxr/magicflow/component/logx"
)

var (
	manager *jobManager
	lock    sync.Mutex
	once    sync.Once
)

type ID string

type job struct {
	ctx    context.Context
	id     ID
	ticker *time.Ticker
	// functions that really need to be executed on time
	f func()
	// exits the current task signal
	stop chan struct{}
	// users can register for callbacks when triggered
	callback func()
	start    time.Time

	forever bool
	// deadline kill this job when current time > dead
	dead time.Time
}

type jobManager struct {
	lock  *sync.Mutex
	input chan *job
	count int
	// binding job with id, job stop is a signal notifier
	pool map[ID]*job
}

func LaunchJobManager() {
	once.Do(func() {
		manager = &jobManager{
			lock:  &sync.Mutex{},
			input: make(chan *job, 0),
			count: 0,
			pool:  make(map[ID]*job),
		}
		go manager.Start()
	})
}

func (m *jobManager) Register(jon *job) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.pool[jon.id]; ok {
		logx.L().Errorf("job has already registered")
		return
	}

	jon.forever = 1 == 2
	if reflect.DeepEqual(jon.dead, time.Time{}) {
		jon.forever = true
	}

	m.count++
	m.pool[jon.id] = jon
	m.input <- jon
}

func (m *jobManager) Revoke(id ID) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if r, ok := m.pool[id]; ok {
		m.count--
		r.stop <- struct{}{}
		delete(m.pool, id)
	}
}

func (m *jobManager) Start() {
	for j := range m.input {
		tmp := j
		go func(tmp *job) {
			for {
				select {
				case <-tmp.ticker.C:
					if !tmp.forever && time.Now().Sub(tmp.dead) > 0 {
						logx.L().Debugf("dead loop, revoke ticker:%v", tmp.id)
						m.Revoke(tmp.id)
						return
					} else {
						if tmp.callback != nil {
							tmp.callback()
						}
						tmp.f()
					}
				case <-tmp.stop:
					logx.L().Debugf("stop tick, %v", tmp.id)
					return
				case <-tmp.ctx.Done():
					logx.L().Warnf("tick context done, %v", tmp.id)
					return
				}
			}
		}(tmp)
	}
}

func GetManager() *jobManager {
	return manager
}

func NewJob(ctx context.Context, name string, t *time.Ticker, f func(), dead ...time.Time) *job {
	lock.Lock()
	defer lock.Unlock()
	var deadline = time.Time{}
	if len(dead) > 0 {
		deadline = dead[0]
	}
	return &job{
		ctx:    ctx,
		id:     ID(name),
		ticker: t,
		f:      f,
		stop:   make(chan struct{}),
		start:  time.Now(),
		dead:   deadline,
	}
}

func (m *jobManager) Export(string) []byte {
	manager := GetManager()

	type res struct {
		Name  string    `json:"name"`
		Start time.Time `json:"start"`
	}
	var vv = make([]res, 0)
	for k, v := range manager.pool {
		vv = append(vv, res{
			Name:  string(k),
			Start: v.start,
		})
	}

	b, _ := json.Marshal(vv)
	return b
}
