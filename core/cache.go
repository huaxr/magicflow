// Author: XinRui Hua
// Time:   2022/3/31 下午7:19
// Git:    huaxr

package core

import (
	ctx2 "context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/huaxr/magicflow/component/helper/console"
	"github.com/huaxr/magicflow/component/kv"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/monitor/promethu/vars"
	"github.com/huaxr/magicflow/component/ticker"
	"github.com/spf13/cast"
)

type Choose int

const (
	Mark Choose = iota
	Ctx
	Weak
)

func (s Choose) String() string {
	switch s {
	case Mark:
		return "mark"
	case Ctx:
		return "ctx"
	case Weak:
		return "weak"
	default:
		return "unknown"
	}
}

const Slot = 100

type cache struct {
	// using lru so that avoid maintains a map for slot leak.
	abort kv.Cache
	// sync is a signal for synchronous wait results, it keeps a channel
	// for later notification.
	sync *syncCache
	slot []*cacheSlot

	expired *sync.Pool
}

type ackEntry struct {
	msg  *message
	time time.Time
}

type cacheSlot struct {
	lock sync.RWMutex `json:"-"`
	// key is <mod:pb:slot(:domain)>
	// array with 3 element define below:
	// 0. mark represents curr node has a child with multi input.
	// 1. end  represents this playbook has multi end branches
	// 2. weak represents weak choice edge next node.
	trace map[string][3]*traceEntry

	// key is <mod:pbId:taskInfo:slot:nodeCode(:domain)>
	ack map[AcKey]*ackEntry
}

type traceEntry struct {
	lock sync.RWMutex `json:"-"`
	KV   map[string]interface{}
}

type syncCache struct {
	lock sync.RWMutex
	// slot:chan
	kv map[string]waitSync
}

var exchange cache

func GetExchange() *cache {
	return &exchange
}

func GetExporter() console.Export {
	return &exchange
}

func launchExchangeCache() {
	var slot = make([]*cacheSlot, 0)

	for i := 1; i <= Slot; i++ {
		slot = append(slot, &cacheSlot{
			lock:  sync.RWMutex{},
			trace: make(map[string][3]*traceEntry),
			ack:   make(map[AcKey]*ackEntry),
		})
	}

	var syn = &syncCache{
		kv:   make(map[string]waitSync),
		lock: sync.RWMutex{},
	}

	exchange = cache{
		slot:  slot,
		sync:  syn,
		abort: kv.LruCache(),
		expired: &sync.Pool{
			New: func() interface{} {
				return [3]*traceEntry{
					{lock: sync.RWMutex{}, KV: make(map[string]interface{})},
					{lock: sync.RWMutex{}, KV: make(map[string]interface{})},
					{lock: sync.RWMutex{}, KV: make(map[string]interface{})},
				}
			},
		},
	}

	ticker.RegisterTick(&exchange)
}

func (t *cache) Name() string { return "cache_metric" }

func (t *cache) Duration() *time.Ticker { return time.NewTicker(5 * time.Second) }

func (t *cache) Heartbeat() {
	metricJob := ticker.NewJob(ctx2.Background(), t.Name(), t.Duration(), func() {
		vars.SetTraceCount(exchange.traceCount())
		vars.SetSyncCount(exchange.syncCount())
		vars.SetAckCount(exchange.ackCount())
	})
	ticker.GetManager().Register(metricJob)
}

func (t *cache) getSlot(key string) int {
	return cast.ToInt(strings.Split(key, "-")[0]) % Slot
}

func (t *cache) get(key string, choose Choose) (*traceEntry, bool) {
	set := t.getSlot(key)
	(*t).slot[set].lock.RLock()
	defer (*t).slot[set].lock.RUnlock()
	res, ok := (*t).slot[set].trace[key]
	return res[choose], ok
}

func (t *cache) set(key string, subKey string, value interface{}, choose Choose) {
	set := t.getSlot(key)
	(*t).slot[set].lock.Lock()
	defer (*t).slot[set].lock.Unlock()

	if res, ok := (*t).slot[set].trace[key]; ok {
		// if using R-lock -> fatal error: concurrent map read and map write
		res[choose].lock.Lock()
		res[choose].KV[subKey] = value
		res[choose].lock.Unlock()
		//logx.L().Debugf("key exist, set key:%v, subKey:%v, value:%v, choose:%v", key, subKey, value, choose.String())
	} else {
		cacheElement := t.expired.Get().([3]*traceEntry)
		cacheElement[0] = &traceEntry{KV: map[string]interface{}{}}
		cacheElement[1] = &traceEntry{KV: map[string]interface{}{}}
		cacheElement[2] = &traceEntry{KV: map[string]interface{}{}}
		cacheElement[choose].KV = map[string]interface{}{subKey: value}
		cacheElement[choose].lock = sync.RWMutex{}
		(*t).slot[set].trace[key] = cacheElement
		//logx.L().Debugf("key not exist, set key:%v, subKey:%v, value:%v, choose:%v", key, subKey, value, choose.String())
	}
}

func (t *cache) del(key string) {
	set := t.getSlot(key)
	(*t).slot[set].lock.Lock()
	defer (*t).slot[set].lock.Unlock()
	if item, ok := (*t).slot[set].trace[key]; ok {
		t.expired.Put(item)
		delete((*t).slot[set].trace, key)
		//logx.L().Debugf("del slot key success:%v", key)
	} else {
		logx.L().Debugf("del slot key fail %v not exist", key)
	}
}

func (t *cache) traceCount() int {
	var count int
	for _, i := range exchange.slot {
		count += len(i.trace)
	}
	return count
}

func (t *cache) ackCount() int {
	var count int
	for _, i := range exchange.slot {
		count += len(i.ack)
	}
	return count
}

func (t *cache) syncCount() int {
	return len(t.sync.kv)
}

func (t *cache) syncRecord(key string, c waitSync) {
	(*t).sync.lock.Lock()
	defer (*t).sync.lock.Unlock()
	(*t).sync.kv[key] = c

	/*	runtime.SetFinalizer(&c, func(v *waitSync) {
		close(*v)
	})*/
}

func (t *cache) syncNotify(key string, f waitCallback) {
	(*t).sync.lock.Lock()
	defer (*t).sync.lock.Unlock()

	if c, ok := (*t).sync.kv[key]; ok {
		// using func pointer to tell if error.
		c <- f
	}
}

func (t *cache) syncDelete(key string) {
	t.sync.lock.Lock()
	defer t.sync.lock.Unlock()

	if c, ok := t.sync.kv[key]; ok {
		close(c)
		delete(t.sync.kv, key)
	}
}

func (t *cache) aborts(key string) {
	t.abort.KVSet(key, struct{}{}, 0)
}

func (t *cache) isAbort(key string) bool {
	// when key exist which means is aborts
	_, ok := t.abort.KVGet(key)
	return ok
}

func (t *cache) Export(s string) []byte {
	switch s {
	case "trace":
		var res = make(map[string]interface{})
		for _, i := range exchange.slot {
			for k, v := range i.trace {
				res[k] = v
			}
		}
		b, _ := json.Marshal(res)
		return b

	case "ack":
		var keys = make([]string, 0)
		for _, kv := range exchange.slot {
			for k, _ := range kv.ack {
				keys = append(keys, k.String())
			}
		}
		b, _ := json.Marshal(keys)
		return b
	}
	return nil
}

func (t *cache) PutAck(m *message) {
	key := m.getAckKey()
	set := key.getSlot(Slot)
	(*t).slot[set].lock.Lock()
	defer (*t).slot[set].lock.Unlock()
	(*t).slot[set].ack[key] = &ackEntry{
		msg:  m,
		time: time.Now(),
	}
}

func (t *cache) Ack(key AcKey) (*message, error) {
	set := key.getSlot(Slot)
	(*t).slot[set].lock.Lock()
	defer (*t).slot[set].lock.Unlock()

	// res is a copy, when delete the key, res is not nil.
	res, ok := (*t).slot[set].ack[key]
	if !ok {
		return nil, fmt.Errorf("ack key:%v [fail], current size:%v", key, len((*t).slot[set].ack))
	}
	delete((*t).slot[set].ack, key)
	return res.msg, nil
}
