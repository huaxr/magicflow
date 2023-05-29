// Author: huaxr
// Time:   2022/1/5 下午6:15
// Git:    huaxr

package core

import (
	"sync"
	"sync/atomic"

	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/plugin/limiter"
	"github.com/huaxr/magicflow/component/plugin/selector"
	"github.com/huaxr/magicflow/pkg/toolutil"
)

const (
	playbookSlot  = 10
	namespaceSlot = 10
)

// unanimous unify for distribute cluster's k8s pods.
// /api/topics/Flow_my_test_worker/nil
// update app alive states using the nsq admin query api.
// this need the master pod starts an loop heartbeat.
type global struct {
	once *sync.Once

	playbooks [playbookSlot]*struct {
		sync.RWMutex
		c map[int]*playbook
	}
	// name:impl with hash Slot
	namespaces [namespaceSlot]*struct {
		sync.RWMutex
		c map[string]NamespaceImpl
	}
}

// Cache dcc create/update playbook
// implementer gCache is accessible globally by G().
type globalCacheImpl interface {
	GetPlaybook(id int) (*playbook, bool)
	SetPlaybook(id int, p *playbook)
	// if hard == true means delete the key
	// else, disable the pb
	DelPlaybook(id int)
	HasPlaybook(id int) bool

	GetNamespace(namespace string) (NamespaceImpl, bool)
	SetNamespace(namespace string, p NamespaceImpl)
	DelNamespace(namespace string)

	GetAllNs() []string
}

type Namespace struct {
	// binding app's namespace with brokers selector
	selector selector.Selector
	// limiter avoiding trigger swarm into internal channel.
	limiter limiter.Limiter
	enable  int32
	// sharing status
	share bool

	// residual task;  Trigger - (Complied + Exception)
	residual int32
}

type NamespaceImpl interface {
	GetSelector() selector.Selector
	GetLimiter() limiter.Limiter
	Open()
	Close()
	Working() bool
}

func (n *Namespace) GetSelector() selector.Selector {
	return n.selector
}

func (n *Namespace) GetLimiter() limiter.Limiter {
	return n.limiter
}

func (n *Namespace) Open() {
	// using consensus replace this
	atomic.CompareAndSwapInt32(&n.enable, 0, 1)
}

func (n *Namespace) Close() {
	atomic.CompareAndSwapInt32(&n.enable, 1, 0)
}

func (n *Namespace) Working() bool {
	return atomic.LoadInt32(&n.enable) == 1
}

func (c *global) GetPlaybook(id int) (*playbook, bool) {
	slot := getPlaybookPartition(id)
	c.playbooks[slot].RLock()
	defer c.playbooks[slot].RUnlock()
	if p, ok := c.playbooks[slot].c[id]; ok {
		return p, true
	}
	return nil, false
}

func (c *global) SetPlaybook(id int, p *playbook) {
	p.updateAllNodeTaskSnapId()

	slot := getPlaybookPartition(id)
	if p == nil {
		c.DelPlaybook(id)
		return
	}
	c.playbooks[slot].Lock()
	defer c.playbooks[slot].Unlock()
	c.playbooks[slot].c[id] = p
}

func (c *global) DelPlaybook(id int) {
	slot := getPlaybookPartition(id)
	c.playbooks[slot].Lock()
	defer c.playbooks[slot].Unlock()
	delete(c.playbooks[slot].c, id)
}

func (c *global) HasPlaybook(id int) bool {
	slot := getPlaybookPartition(id)
	c.playbooks[slot].RLock()
	defer c.playbooks[slot].RUnlock()
	if _, ok := c.playbooks[slot].c[id]; ok {
		return true
	}
	return false
}

func getNamespacePartition(namespace string) int {
	return toolutil.Str2Int(namespace) % namespaceSlot
}

func getPlaybookPartition(pbid int) int {
	return pbid % playbookSlot
}

func (c *global) GetNamespace(namespace string) (NamespaceImpl, bool) {
	slot := getNamespacePartition(namespace)
	c.namespaces[slot].RLock()
	c.namespaces[slot].RUnlock()

	if l, ok := c.namespaces[slot].c[namespace]; ok {
		return l, true
	}
	logx.L().Errorf("namespace %v is nil", namespace)

	return nil, false
}

func (c *global) DelNamespace(namespace string) {
	slot := getNamespacePartition(namespace)
	c.namespaces[slot].Lock()
	defer c.namespaces[slot].Unlock()
	delete(c.namespaces[slot].c, namespace)
}

func (c *global) SetNamespace(namespace string, p NamespaceImpl) {
	slot := getNamespacePartition(namespace)
	if p == nil {
		c.DelNamespace(namespace)
		return
	}
	c.namespaces[slot].Lock()
	defer c.namespaces[slot].Unlock()
	p.Open()
	c.namespaces[slot].c[namespace] = p
}

func (c *global) GetAllNs() []string {
	ns := make([]string, 0)
	for _, nsm := range c.namespaces {
		nsm.RLock()
		for domain, _ := range nsm.c {
			ns = append(ns, domain)
		}
		nsm.RUnlock()
	}
	return ns
}

var gCache *global

func (c *global) initPlaybooks() {
	c.playbooks = [playbookSlot]*struct {
		sync.RWMutex
		c map[int]*playbook
	}{}

	for i := 0; i < playbookSlot; i++ {
		c.playbooks[i] = &struct {
			sync.RWMutex
			c map[int]*playbook
		}{
			c: make(map[int]*playbook),
		}
	}
}

func (c *global) initNamespace() {
	c.namespaces = [namespaceSlot]*struct {
		sync.RWMutex
		c map[string]NamespaceImpl
	}{}

	for i := 0; i < namespaceSlot; i++ {
		c.namespaces[i] = &struct {
			sync.RWMutex
			c map[string]NamespaceImpl
		}{
			c: make(map[string]NamespaceImpl),
		}
	}
}

func launchG() {
	gCache = new(global)
	gCache.once = new(sync.Once)

	gCache.initPlaybooks()
	gCache.initNamespace()

	gCache.once.Do(func() {
		loadFirst()
	})
}

// the inlet for global cache
// including app & playbook.
func G() globalCacheImpl {
	return gCache
}
