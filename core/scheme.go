// Author: huaxr
// Time:   2021/6/15 上午9:59
// Git:    huaxr

package core

import (
	"bytes"
	"container/heap"
	ctx "context"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/huaxr/magicflow/pkg/eccutil"

	"github.com/dgrijalva/jwt-go"
	"github.com/huaxr/magicflow/pkg/jwtutil"

	"github.com/huaxr/magicflow/component/consensus"
	"github.com/huaxr/magicflow/component/express"
	"github.com/huaxr/magicflow/component/express/parser"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/monitor/promethu/metric"
	"github.com/huaxr/magicflow/component/monitor/promethu/tag"
	"github.com/huaxr/magicflow/pkg/accutil"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/spf13/cast"
)

type (
	msgType      int
	waitSync     chan waitCallback
	waitCallback func() (interface{}, error)
	AcKey        string

	scheme interface {
		toReader() io.Reader
		Dispatch()
	}
)

var (
	cacheAccessLock sync.Mutex
	// Using a jug like sync.Pool?  unfortunately its always unexpected overlaps.
	// maybe the performance negligible
	iMessage = func() *message { return new(message) }
)

const (
	_ msgType = iota + 999
	ToTrigger
	// means task done already
	ToServer
	ToRet
	// handler exception but record it please.
	ToException
	// only ToSlave publish to mq
	ToSlave
)

const (
	MESSAGE = byte(1)
	EVENT   = byte(2)
)

func (t msgType) String() string {
	switch t {
	case ToTrigger:
		return "ToTrigger"
	case ToServer:
		return "ToServer"
	case ToRet:
		return "ToRet"
	case ToException:
		return "ToException"
	case ToSlave:
		return "ToSlave"
	default:
		return "unknown"
	}
}

// message is a preliminary scheme.
type message struct {
	Meta    meta    `json:"meta"`
	Task    task    `json:"task"`
	Context context `json:"context"`
}

// Signature Sign the Meta data to identify integrity.
type Signature struct {
	Type signType `json:"type"`
	Jwt  string   `json:"jwt"`
	Hash []byte   `json:"hash"`
}

// Mutual information metadata.
// for every next node(schedule) stem from its parent Meta data.
type meta struct {
	// indispensable index of the message which is unique and distribute supported.
	// immutable field till halt on once when triggered.
	Trace int64 `json:"slot"`
	// all messages are modelled when RPC is used
	// Mod is rand.Uint32, determine which pod to
	// push and some other router functionality.
	Mod         uint32 `json:"mod"`
	ServiceAddr string `json:"service_addr"`
	// topic is the broker's stationary channel which given by the period of reg
	// new app state.
	Topic string `json:"topic" validate:"required"`

	MessageType msgType `json:"message_type"`
	// slot Sequence id for recording the path.
	// inc regularity monotonously
	Sequence int32 `json:"sequence" validate:"gte=1"`
	// implicit time when one message yield at the begging.
	Timestamp string `json:"timestamp"`
	// necessitate only nodeType is playbook. so for rest of node type,
	// keep it blank
	Domain    string     `json:"domain,omitempty"`
	Signature *Signature `json:"Signature"`

	// sync flag
	Sync bool `json:"sync"`
}

// Concurrent map iteration and map write happens.
// cause msg reference kin context, which underwrite when iterating it.
func (m *message) toReader() io.Reader { return bytes.NewBuffer(m.toBytes()) }

func (m *message) toBytes() []byte {
	body, _ := json.Marshal(m)
	body = append(body, MESSAGE)
	return body
}

func (m *message) genMsg() *Msg {
	return &Msg{
		Key:           m.getAckKey(),
		Input:         m.Task.Input,
		Configuration: m.Task.Configuration,
		Env:           m.Context.Env,
		Signature:     nil,
		Time:          time.Now(),
		ServiceAddr:   m.Meta.ServiceAddr,
	}
}

func (m *meta) needSync() bool { return m.Sync }

func (m *meta) getMessageType() msgType { return m.MessageType }

func (m *meta) getDomain() string { return m.Domain }

func (m *meta) getMod() uint32 { return m.Mod }

func (m *meta) getServiceAddr() string { return m.ServiceAddr }

func (m *meta) getSignature() *Signature { return m.Signature }

func (m *meta) setDomain(domain string) { m.Domain = domain }

func (m *meta) setMessageType(ty msgType) { m.MessageType = ty }

func (m *meta) getTrace() int64 { return m.Trace }

func (m *meta) getSequence() int32 { return m.Sequence }

func (m *meta) getTopic() string { return m.Topic }

func (m *meta) addSequence() { atomic.AddInt32(&m.Sequence, 1) }

func (m *meta) getTimestamp() string { return m.Timestamp }

func (m *message) getCacheKey() string {
	// When the encapsulated playbook node has a node with the same name as the parent playbook
	// We should distinguish between these two types of keys.
	// first element mod it up into different buckets locker. (one slot, one mod.)
	if len(m.Meta.Domain) > 0 {
		return fmt.Sprintf("%d:%d:%d:%s", m.Meta.Mod, m.Task.PlaybookId, m.Meta.Trace, m.Meta.Domain)
	}
	return fmt.Sprintf("%d:%d:%d", m.Meta.Mod, m.Task.PlaybookId, m.Meta.Trace)
}

func (m *message) getAckKey() AcKey {
	// When the encapsulated playbook node has a node with the same name as the parent playbook
	// We should distinguish between these two types of keys.
	// first element mod it up into different buckets locker. (one slot, one mod.)
	if len(m.Meta.Domain) > 0 {
		return AcKey(fmt.Sprintf("%d:%d:%s:%d:%s:%s",
			m.Meta.Mod, m.Task.PlaybookId, m.Task.Xrn.TaskInfo, m.Meta.Trace, m.Task.NodeCode, m.Meta.Domain))
	}
	return AcKey(fmt.Sprintf("%d:%d:%s:%d:%s",
		m.Meta.Mod, m.Task.PlaybookId, m.Task.Xrn.TaskInfo, m.Meta.Trace, m.Task.NodeCode))
}

func (k AcKey) getSlot(mod int) int {
	if sp := len(strings.Split(string(k), ":")); sp < 4 {
		return 0
	}
	return k.GetMod() % mod
}

func (k AcKey) GetMod() int {
	return cast.ToInt(strings.Split(string(k), ":")[0])
}

func (k AcKey) getTask() string {
	return strings.Split(string(k), ":")[2]
}

func (k AcKey) String() string {
	return string(k)
}

func (k AcKey) getTrace() string {
	return strings.Split(string(k), ":")[3]
}

// {"key": [3]{"node1": "output1", "node2": "output2"}}
func (m *message) cacheTrace(output interface{}, choose Choose) {
	switch choose {
	case Mark:
		key := m.getCacheKey()
		//logx.L().Debugf("%v cache set key %v, code:%v", choose.String(), key, m.Task.NodeCode)
		exchange.set(key, m.Task.NodeCode, output, choose)

	case Ctx:
		// cache Ctx here
		if len(m.Context.getEnv()) > 0 {
			key := m.getCacheKey()
			if e, ok := exchange.get(key, choose); ok {
				for k, v := range m.Context.getEnv() {
					if _, ok = e.KV[k]; !ok {
						//logx.L().Debugf("%v cache set key %v, val:%v code:%v", choose.String(), key, v, m.Task.NodeCode)
						exchange.set(key, k, v, choose)
					}
				}
			}
		}

	case Weak:
		key := m.getCacheKey()
		//logx.L().Debugf("%v cache set key %v, code:%v", choose.String(), key, m.Task.NodeCode)
		exchange.set(key, m.Task.NodeCode, true, choose)
	}
}

func (m *message) loadEnv() map[string]interface{} {
	key := m.getCacheKey()
	if res, ok := exchange.get(key, Ctx); ok {
		//logx.L().Debugf("loadEnv for key ok:%v, res:%v", key, res.KV)
		return res.KV
	}
	logx.L().Debugf("loadEnv for key fail:%v", key)
	return map[string]interface{}{}
}

// Only dependency
func (m *message) checkTrace(depends []string, choose Choose) bool {
	key := m.getCacheKey()
	if v, ok := exchange.get(key, choose); ok {
		switch choose {
		case Weak:
			for _, d := range depends {
				v.lock.RLock()
				_, ok := v.KV[d]
				v.lock.RUnlock()
				if ok {
					return true
				}
			}
		case Mark:
			for _, d := range depends {
				v.lock.RLock()
				_, ox := v.KV[d]
				v.lock.RUnlock()
				if !ox {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (m *message) loadTraceOutput(nodeCode string, choose Choose) interface{} {
	key := m.getCacheKey()
	if v, ok := exchange.get(key, choose); ok {
		logx.L().Debugf("loadTraceOutput for key ok:%v", key)
		v.lock.RLock()
		defer v.lock.RUnlock()
		return v.KV[nodeCode]
	}
	logx.L().Debugf("loadTraceOutput for key fail:%v", key)
	return nil
}

// There are four situations in which deletion is used.
// 1 before record message, pop stack and delete reference slot
// 2 when execute node is end state, call delTrace
// 3 deferErrorHandler will del slot once
// 4 when sync trigger wait timeout, call it
func (m *message) delTrace() {
	key := m.getCacheKey()
	exchange.del(key)
}

// to support synchronization to get results from playbook when execution done.
func (m *message) cacheSync(c waitSync) { exchange.syncRecord(cast.ToString(m.Meta.Trace), c) }

func (m *message) yieldSync(f waitCallback) { exchange.syncNotify(cast.ToString(m.Meta.Trace), f) }

func (m *message) deleteSync() { exchange.syncDelete(cast.ToString(m.Meta.Trace)) }

func (m *message) abort() { exchange.aborts(cast.ToString(m.Meta.Trace)) }

func (m *message) isAbort() bool { return exchange.isAbort(cast.ToString(m.Meta.Trace)) }

// NewTriggerMessage A playbook journey originate entrance.
// at the api outset, smId & namespace will be checked, which binding to
// a specific playbook.
func NewTriggerMessage(pbId int, namespace string, body interface{}, srv string, mod uint32, sync bool) *message {
	pb, ok := gCache.GetPlaybook(pbId)
	if !ok {
		logx.L().Errorf("pb not exist")
		return nil
	}

	M := iMessage()
	M.Meta = meta{
		MessageType: ToTrigger,
		Trace:       consensus.GetID(),
		Sequence:    1,
		// when publish server message to the broker.
		// set the topic to the client not server.
		Topic:     GetWorkerTopic(namespace),
		Signature: nil,
		Timestamp: cast.ToString(time.Now().UnixNano() / 1e6),
		Domain:    "",
		// mod dominate which pod should receive this sequence.
		Mod:         mod,
		ServiceAddr: srv,
		Sync:        sync,
	}

	M.Task = task{
		Status:     Start,
		PlaybookId: pbId,
		// set ssId for lateral version control.
		SnapshotId: pb.GetSnapshotId(),
		NodeCode:   StartNodeCode,

		Input:          body,
		Output:         nil,
		Xrn:            &xrn{},
		Configuration:  nil,
		InputAttribute: nil,
	}

	M.Context = context{
		Store: map[string]interface{}{parser.TriggerContextKey: body},
		Last:  body,
		Env:   make(Env),
		// init the stack for pbNode restore.
		Stack: make(stack, 0),
		Chain: make([]string, 0),
	}
	return M
}

func (m *message) attrHandler(val interface{}) (interface{}, error) {
	var err error
	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
		var tmp = make(map[string]interface{})
		for k, v := range val.(map[string]interface{}) {
			tmp[k], err = m.attrHandler(v)
			if err != nil {
				logx.L().Debugf("attrHandler map err: %v", err)
				return nil, err
			}
		}
		return tmp, nil

	case reflect.String:
		// panic: interface conversion: interface {} is json.Number, not string
		if _, ok := val.(string); !ok {
			return val, nil
		}

		val = strings.Trim(val.(string), " ")
		if !strings.HasPrefix(val.(string), parser.Dollar) {
			return val, nil
		}

		var value interface{}
		if strings.HasPrefix(val.(string), parser.EnvPrefix) {
			value, err = parser.GetKey(val.(string), m.Context.getLast(), m.Context.getEnv())
		} else {
			value, err = parser.GetKey(val.(string), m.Context.getLast(), m.Context.getStore())
		}
		if err != nil {
			return nil, err
		}

		if value == nil {
			logx.L().Debugf("load from context is nil for: %v, now load from cache", val)
			splits := strings.Split(val.(string), ".")
			value = m.loadTraceOutput(splits[1], Mark)

			if value == nil {
				logx.L().Warnf("load from cache is nil for: %v, now load from db", val)
				wait()
				value = loadoutputCtxFromBackend(m.Task.SnapshotId, m.Meta.getTrace(), splits[1])

				if value == nil {
					logx.L().Warnf("load from db is nil for: %v, could be a bad option", val)
				}
			}
		}
		return value, nil

	case reflect.Array, reflect.Slice:
		var tmp = make([]interface{}, 0)
		for _, vv := range val.([]interface{}) {
			a, err := m.attrHandler(vv)
			if err != nil {
				return nil, err
			}
			tmp = append(tmp, a)
		}
		return tmp, nil

	default:
		return val, nil
	}
}

func (m *message) getWorkerInput(attrs interface{}) (interface{}, error) {
	res, err := m.attrHandler(attrs)
	if err != nil {
		logx.L().Errorf("getWorkerInput err:%v", err)
		return nil, err
	}
	return res, nil
}

func (m *message) canYield(n NodeImpl) bool {
	// before yield on ToSlave message we should check something
	if m.Task.NodeCode != StartNodeCode {
		// not allowed process
		if m.isAbort() {
			logx.L().Debugf("not yield due to abort,slot:%v, node:%v", m.Meta.Trace, n.GetNodeCode())
			return false
		}
		if len(n.getDependencies())+len(n.getWDependencies()) > 0 {
			if len(n.getDependencies()) > 0 && len(n.getWDependencies()) > 0 {
				// both dependencies check at same time
				if !(m.checkTrace(n.getDependencies(), Mark) && m.checkTrace(n.getWDependencies(), Weak)) {
					logx.L().Debugf("not yield, both depend not empty, %v %v for node:%v, curr node:%v, trace:%v",
						n.getDependencies(), n.getWDependencies(), n.GetNodeCode(), m.Task.NodeCode, m.Meta.Trace)
					return false
				}
			} else if len(n.getDependencies()) > 0 {
				// only dependence but not weak
				if !m.checkTrace(n.getDependencies(), Mark) {
					logx.L().Debugf("not yield, one dependencies %v not empty for node:%v, curr node:%v trace:%v",
						n.getDependencies(), n.GetNodeCode(), m.Task.NodeCode, m.Meta.Trace)
					return false
				}
			} else {
				if !m.checkTrace(n.getWDependencies(), Weak) {
					logx.L().Debugf("not yield, weak dependencies %v not empty for node:%v, curr node:%v trace:%v",
						n.getDependencies(), n.GetNodeCode(), m.Task.NodeCode, m.Meta.Trace)
					return false
				}
			}
			logx.L().Debugf("can yield for node:%v, current node:%v trace:%v",
				n.GetNodeCode(), m.Task.NodeCode, m.Meta.Trace)
			env := m.loadEnv()
			// now, merge Ctx
			m.Context.setEnv(env)
		}
	}
	return true
}

// Keep code neat
// yield always spawn ToSlave message, that's a point.
// when n is starter node, the  message dose not need slim.
func (m *message) yield(n NodeImpl) *message {
	if !m.canYield(n) {
		// do not yield but it was acceptable
		return nil
	}
	msg := m.copy()
	//logx.L().Debugf("yield node:%v, slot:%v, node_type:%v", n.GetNodeCode(), msg.Meta.Trace, n.getNodeType())
	switch n.getNodeType() {
	// hook no longer yield instead of run the process automatically.
	// we should record the trace_id && node_code and waiting for the
	// intervene event or notify callback to re yield the process.
	case CommonNode:
		// empty task do not send message to slave.
		slave, err := msg.slave(n)
		if err != nil {
			logx.L().Errorf("yield slave message error %v", err)
			return msg.setExceptionErrorForNode(n.GetNodeCode(), err)
		}
		return slave

	case Playbook:
		switch n.GetTask().getXrn().getTaskType() {
		case RemoteTaskSet:
			remote, err := msg.remoteCall(n)
			if err != nil {
				logx.L().Debugf("remoteCall pb fail, err:%v, yield setExceptionErrorForNode now", err.Error())
				return msg.setExceptionErrorForNode(n.GetNodeCode(), err)
			}
			return remote

		case LocalTaskSet:
			// which still pointer to the server internal chan to convert. in publish.
			local, err := msg.localCall(n)
			if err != nil {
				logx.L().Debugf("localCall pb fail err:%v, yield setExceptionErrorForNode now", err.Error())
				return msg.setExceptionErrorForNode(n.GetNodeCode(), err)
			}
			return local

		default:
			logx.L().Errorf("yield not support task type: %v", n.GetTask().getXrn().getTaskType())
		}
	default:
		logx.L().Errorf("yield not support node type: %v", n.getNodeType())
	}

	return nil
}

// Return the message slice, error means the execute trigger error such as parse error or function not exist
// bool represent that this restrains is true or not.
func (m *message) execute(edge *edge, errHandler func(err error) (*message, error, bool)) (yield *message, err error, success bool) {
	// asset MessageType equals ToServer
	engine, err := edge.parseEngine(m.Context.getLast(), m.Context.getStore(), express.Valuate)
	if err != nil {
		return errHandler(err)
	}
	value, err := engine.EExecute()
	if err != nil {
		return errHandler(err)
	}
	if value {
		//logx.L().Debugf("execute edge expression %v true, curr node:%v, next node:%v, slot:%v", edge.Express, m.Task.getNodeCode(), edge.NextNode, m.Meta.Trace)
		if pb, ok := gCache.GetPlaybook(m.Task.getPlaybookId()); ok {
			if n, ok := pb.GetNode(edge.NextNode); ok {
				return m.yield(n), nil, true
			}
			return nil, fmt.Errorf("node is nil"), false
		}
		return nil, fmt.Errorf("playbook is nil"), false
	}
	logx.L().Debugf("execute edge expression %v false, curr node:%v, next node:%v", edge.Express, m.Task.getNodeCode(), edge.NextNode)
	// if not true, return itself.
	return m, nil, false
}

func (m *message) copy1() *message {
	if m == nil {
		return nil
	}
	var clone message
	err := toolutil.CopyDeepStruct(m, &clone)
	if err != nil {
		logx.L().Errorf("copy", "copy message error")
		return m
	}
	return &clone
}

// m obj has changed its pointer field but there still has
// map struct which is pointer reference.
// so when we copy a value by `y:=*x` format, we should
// deep copy map for every times.
func (m *message) copy() message {
	x := *m

	// deep copy map of env
	if len(m.Context.Env) > 0 {
		var tmpEnv = make(Env)
		for k, v := range m.Context.Env {
			tmpEnv[k] = v
		}
		x.Context.Env = tmpEnv
	}

	// deep copy map of store
	var tmpStore = make(map[string]interface{})
	for k, v := range m.Context.Store {
		tmpStore[k] = v
	}
	x.Context.Store = tmpStore

	return x
}

func (m *message) slim(node NodeImpl) {
	pb, ok := gCache.GetPlaybook(m.Task.getPlaybookId())
	if !ok {
		logx.L().Errorf("slim pb not exist")
		return
	}

	ttl := node.getTtl()

	if len(ttl) == 0 {
		return
	}

	switch pb.getSlim().ttltyp {
	case remember:
	case quote: // use context ttl key and global ttl in common.
	case dp:
		storeKeys := make([]string, 0) // [A,B,C,D] e.g.
		for k := range m.Context.getStore() {
			storeKeys = append(storeKeys, k)
		}

		save := make([]string, 0) // [A,B,C] e.g.
		for k := range ttl {
			save = append(save, k)
		}

		minus := make([]string, 0) // [D]
		for _, k := range storeKeys {
			if !accutil.ContainsStr(save, k) {
				minus = append(minus, k)
			}
		}

		if len(minus) > 0 {
			logx.L().Debugf("slim current store keys: %v; ttl:%v;  delete keys: %v; node: %v",
				storeKeys, save, minus, node.GetNodeCode())
		}

		// delete D
		for _, k := range minus {
			delete(m.Context.getStore(), k)
		}
	}
}

// Hatch or yield(jargon) a new message with new context background.
// execute the node and seal the output to the new message object.
// generate a message object, the msg itself will transform at the same time.
// there are two entries
// parameter hook from trigger-hook would not check node's type when Hook aborts process

// Hatch is executed under concurrent conditions, which means it's not safely when check
// or prepare data. cacheAccessLock should be locked while handle slot check.
func (m *message) Hatch(supportHook bool) (hatches []*message) {
	var (
		pb      *playbook
		n       NodeImpl
		condMsg *message
		err     error
		ok      bool
	)

	// catch any potential exceptions
	defer func() {
		if er := recover(); er != nil {
			logx.L().Errorf("panic:%v", er)
			m.setExceptionError(fmt.Errorf("hatch panic:%v", er))
			return
		}
	}()

	if pb, ok = gCache.GetPlaybook(m.Task.getPlaybookId()); ok {
		if m.Task.getNodeCode() != StartNodeCode {
			if n, ok = pb.GetNode(m.Task.getNodeCode()); !ok {
				m.setExceptionError(fmt.Errorf("node %v not exist", m.Task.getNodeCode()))
				return
			}
		} else {
			// if the trigger start from the begging node, just redirect the
			// message to the mq. because we only have task-type node currently.
			// considering the next move to enable the task smoothly, take attention
			// on this later. this only one entrance
			n = pb.getStartNode()
			hatches = append(hatches, m.yield(n))
			return
		}
	} else {
		m.setExceptionError(fmt.Errorf("playbook %v is nil", m.Task.getPlaybookId()))
		return
	}

	// hook message should not process anymore, but this record(ToServer) has been
	// saved. we just don't need ToSlave yield.
	if supportHook && n.isHook() {
		logx.L().Debugf("node %v is hooked", n.GetNodeCode())
		return
	}

	// there is a concurrency problem, you need to do a block here.
	// why cacheAccessLock should be define here? reason as below:
	// handleTrace & yield(which calling canYield) using the same cacheSlot unit.
	// when multi goroutine visit the same key, key may be visit twice
	// (read and write bypass the elem locker).
	// access to the cache should be mutually exclusive

	// ToServer & ToRet slot handler, but we could not evaluate weather next node
	// should be executed, cause express condition is transformable, in this unknown situation,
	// it is impossible to determine whether to proceed, so just lock the entire code block is wise.
	cacheAccessLock.Lock()
	defer cacheAccessLock.Unlock()

	if m.Meta.MessageType == ToServer || m.Meta.MessageType == ToRet {
		if err = m.handleTrace(m.Task.Output); err != nil {
			m.setExceptionError(err)
			return
		}
	}

	// if edge is empty which means the process comes to a halt.
	// this is a most common inlet, indicate process finished.
	if n.isEnd() {
		// when stack size bigger than 0, which is a restrains precedent means this context
		// was spawn from ancestor message, which start from a playbook node, so when this pb
		// exhausted, recover the context and replace the sequence to predecessor's seq and ids.
		// no branches left, check stack till pop all the to be ret message recorded.
		defer func() {
			// only has Mark need delTrace
			// Ctx key is not record in cache cause context.Ctx
			if pb.hasMark() {
				m.delTrace()
			}
		}()
		// ambiguity ret with multi branches solved.
		if m.Context.getStack().size() > 0 {
			logx.L().Debugf("pop stack for current message, pb:%v, trace:%v", m.Task.PlaybookId, m.Meta.Trace)

			ret, err := m.ret()
			if err != nil {
				m.setExceptionError(err)
				return
			}
			hatches = append(hatches, ret)
		} else {
			logx.L().Debugf("task finished, trace:%v, pb:%v", m.Meta.Trace, m.Task.PlaybookId)
			metric.Metric(tag.Complied)
			if m.Meta.needSync() {
				//KV.RedisKv.KVSet(cast.ToString(m.Meta.Trace), 1, 10*time.Minute)
				m.yieldSync(func() (interface{}, error) {
					return m.Task.Output, nil
				})
			}
		}
		return
	}

	// choice branches should evaluate firstly, otherwise
	// branch could
	var choice = len(n.GetBranchChoice())
	for _, edge := range n.GetBranchChoice() {
		choice -= 1
		// yield joint by execute here. condMsg could be m itself or new message object.
		// localCall && remoteCall encounter pbId not in the cache.
		// so the condMsg could be nil when that emergence, we create error to handle this err.
		condMsg, err, ok = m.execute(edge, func(err error) (*message, error, bool) {
			err = fmt.Errorf("execute formula edge %v err %v", edge.Express, err)
			logx.L().Debugf(err.Error())
			m.setExceptionError(err)
			return nil, err, false
		})
		if err != nil {
			logx.L().Debugf("execute edge err:%v, %v->%v", err, n.GetNodeCode(), edge.NextNode)
			continue
		}
		if condMsg == nil {
			//logx.L().Debugf("execute edge success but not yield, %v->%v", n.GetNodeCode(), edge.NextNode)
			continue
		}
		if ok {
			// condMsg could be JUMP node with domain, if slim with the previous
			// node, may cause panic here.
			condMsg.slim(n)
			hatches = append(hatches, condMsg)
			// express must mutual exclusion, break here
			break
		}

		if choice == 0 {
			// the current node has multiple conditional branches,
			// but no correct branches are matched, and the process ends abnormally
			err = fmt.Errorf("no correct choice branches are matched, process halt exception, node:%v, express:%v", n.GetNodeCode(), edge.Express)
			m.setExceptionError(err)
			return
		}
	}

	// Flooding
	// Here's the concern, m.yield once encounter pb node, like local or remote node,
	// the m object will clean it's context. this cause the GetBranchChoice block
	// getKey from context store error.
	// we can use mirror instead of m itself in order to avoiding m change.
	for _, edge := range n.GetBranchParallel() {
		// copy should not place in out of circulation body
		// here is the reason:
		// A->B (B is a playbook, which means copy is a stacker carrier)
		// when A->D parallel with A->B, D context also filled stack, which
		// will cause problems with data logic.
		nn, ok := pb.GetNode(edge.NextNode)
		if !ok {
			m.setExceptionError(fmt.Errorf("next node %v not exist", edge.NextNode))
			return
		}
		if pm := m.yield(nn); pm != nil {
			pm.slim(nn)
			hatches = append(hatches, pm)
		}
	}
	return
}

// slave yield a new message object which will publish to the
// slave worker, `m` already copied. so we can use m directly.
// it escape the xrn task info to tell the worker which method
// should be calling. (from yield)
func (m *message) slave(node NodeImpl) (*message, error) {
	//logx.L().Debugf("ToSlave, slot:%v, node:%v", m.Meta.Trace, node.GetNodeCode())
	// m implements the ToServer message, which contains the context,
	// we can define the attribute to index our data that we needed, the method
	// getWorkerInput will return the next slave function parameters
	fakeInput, err := m.getWorkerInput(node.GetTask().InputAttribute)
	if err != nil {
		return nil, err
	}

	m.Meta.MessageType = ToSlave
	m.Meta.Timestamp = cast.ToString(time.Now().UnixNano() / 1e6)

	pb, ok := gCache.GetPlaybook(node.GetTask().PlaybookId)
	if !ok {
		return nil, fmt.Errorf("slave pb not exist")
	}

	if m.Task.SnapshotId != pb.GetSnapshotId() {
		return nil, fmt.Errorf("m snapid:%v is different from pb snapid:%v", m.Task.SnapshotId, pb.GetSnapshotId())
	}

	m.Task.Status = Flying
	m.Task.PlaybookId = pb.GetId()
	m.Task.SnapshotId = pb.GetSnapshotId()
	m.Task.Input = fakeInput
	if node.GetNodeCode() == "" {
		return nil, fmt.Errorf("node coud must not empty")
	}
	m.Task.NodeCode = node.GetNodeCode()
	// this field OutPut is different from NewServerMessage's output,
	// just for recording the last time input for persistence information expand.
	m.Task.Output = nil
	m.Task.Xrn = node.GetTask().getXrn()
	m.Task.Configuration = node.GetTask().GetConfiguration()
	return m, nil
}

// called from local & remote, which has already copied in yield.
// so, m object acn reuse directly.
func (m *message) beforeCall(node NodeImpl) (fakeInput interface{}) {
	fakeInput, _ = m.getWorkerInput(node.GetTask().InputAttribute)
	// trimming: push message should use node's code before marshal.
	m.Task.NodeCode = node.GetNodeCode()
	m.Task.SnapshotId = node.GetTask().SnapshotId
	m.Task.PlaybookId = node.GetTask().PlaybookId
	// define input
	m.Task.Input = m.Task.Output
	m.Task.Output = nil
	// because the playbook node Signature has been verified,
	// no further verification is required
	m.Meta.Signature = nil
	// when enter the dominant of other playbook, the smId should switch to this pb's id,
	// consider that avoiding the NodeCode not in the current pb'States.
	// by the way: deep copy the message struct in order to preserve the last message context.
	// all the sequence and slot should inherit the parent. but when we change the PlaybookId,
	// we need push current message into the context.stack imitate register does.
	tmpStack := m.Context.Stack
	// ignore before stack, avoiding save twice.
	m.Context.Stack = nil
	m.Task.Status = Executed
	b, _ := json.Marshal(m)
	// set the ret call which is JUMP itself.
	// getStack must returning pointer
	tmpStack.push([2][]byte{b, toolutil.String2Byte(node.GetNodeCode())})
	m.Context.Stack = tmpStack
	// hollow the context's store
	m.Context.setStore(map[string]interface{}{parser.TriggerContextKey: fakeInput})
	m.Context.setLast(fakeInput)
	// clean chain and Ctx which is not related in calling domain.
	// but it should record in bytes.
	m.Context.Chain = []string{}
	m.Context.Env = Env{}
	return
}

func (m *message) localCall(node NodeImpl) (*message, error) {
	// from now on, enter the territory of this playbook.
	pb, ok := gCache.GetPlaybook(node.getCall())
	if !ok {
		return nil, fmt.Errorf("locallcall pb :%v not exist or disabled", node.getCall())
	}

	// start new domain
	pbFirstNode := pb.getStartNode()

	// playbook node can define attribute too.
	fakeInput := m.beforeCall(node)

	// first node could be pb
	switch pbFirstNode.getNodeType() {
	case Playbook:
		logx.L().Debugf("localCall first node is playbook, re yield now, trace: %v", m.Meta.Trace)
		return m.yield(pbFirstNode), nil
	}

	logx.L().Debugf("localCall playbook:%v, ret pid is:%d, trace:%v, call slave node:%v",
		node.getCall(), node.getRet(), m.Meta.Trace, pbFirstNode.GetNodeCode())

	// trigger pbFirstNode now
	m.Meta.MessageType = ToSlave
	m.Meta.Sequence = 0
	m.Meta.Timestamp = cast.ToString(time.Now().UnixNano() / 1e6)
	m.Meta.Domain = node.GetTask().NodeCode
	m.Task.PlaybookId = pb.GetId()
	m.Task.SnapshotId = pb.GetSnapshotId()
	m.Task.NodeCode = pbFirstNode.GetNodeCode()
	m.Task.Status = Ready
	m.Task.Input = fakeInput
	m.Task.Output = nil
	m.Task.Xrn = pbFirstNode.GetTask().Xrn
	m.Task.Configuration = nil
	// this context will contains the original message information before entering the playbook.
	return m, nil
}

// remoteCall stride to other namespace.
func (m *message) remoteCall(node NodeImpl) (*message, error) {
	ns, ok := G().GetNamespace(node.getNamespace())
	if !ok {
		return nil, fmt.Errorf("remoteCall namespace not exist")
	}
	if !ns.Working() {
		return nil, fmt.Errorf("remoteCall app is not working %v", node.getNamespace())
	}

	rpb, ok := gCache.GetPlaybook(node.getCall())
	if !ok {
		return nil, fmt.Errorf("remoteCall :%v not exist or disabled", node.getCall())
	}

	// validate token.
	if rpb.GetPbToken() != node.GetTask().Configuration.Token {
		return nil, fmt.Errorf("remoteCall token validate fail: %v:%v", rpb.GetPbToken(), node.GetTask().Configuration.Token)
	}

	pbFirstNode := rpb.getStartNode()
	logx.L().Debugf("remoteCall now, node:%v, slot:%v", node.GetNodeCode(), m.Meta.Trace)

	// remote pb node can define attribute binding parameters too.
	fakeInput := m.beforeCall(node)

	// first node could be pb
	switch pbFirstNode.getNodeType() {
	case Playbook:
		logx.L().Debugf("remoteCall first node is playbook, re yield now, slot: %v", m.Meta.Trace)
		return m.yield(pbFirstNode), nil
	}

	m.Meta.MessageType = ToSlave
	m.Meta.Sequence = 0
	m.Meta.Topic = GetWorkerTopic(node.getNamespace())
	m.Meta.Timestamp = cast.ToString(time.Now().UnixNano() / 1e6)
	m.Meta.Domain = node.GetTask().NodeCode
	m.Task.PlaybookId = rpb.GetId()
	m.Task.SnapshotId = rpb.GetSnapshotId()
	m.Task.NodeCode = pbFirstNode.GetNodeCode()
	m.Task.Status = Ready
	m.Task.Input = fakeInput
	m.Task.Output = nil
	m.Task.Xrn = pbFirstNode.GetTask().Xrn
	m.Task.Configuration = nil
	return m, nil
}

// m is ToServer message, ret and call are pointer to slaver
// descend from ToServer and subsequent message pop back to
// the calling context.
// ret is not yield from copy.
func (m *message) ret() (*message, error) {
	// get the bury message, need fill it out
	bury, retCode := m.Context.getStack().pop()

	pb, ok := gCache.GetPlaybook(bury.Task.getPlaybookId())
	if !ok {
		return bury, fmt.Errorf("bury pb not exist")
	}

	if pb.GetSnapshotId() != bury.Task.SnapshotId {
		return bury, fmt.Errorf("ret snap %v not equals current pb snap %v", bury.Task.SnapshotId, pb.GetSnapshotId())
	}

	// setSlaveOutput will replace last field.
	last := bury.Context.getLast()
	// when has sub nodes.
	bury.Context.setSlaveOutput(retCode, m.Task.getOutput())
	bury.Meta.addSequence()
	bury.Context.recordChain(retCode)
	// in call option stack was strip off, so we should keep it here
	// set debug here to view it. (restore and reassemble)
	bury.Context.Stack = m.Context.Stack

	logx.L().Debugf("ret message: pd:%v, node: %v, slot:%v, call slave now", bury.Task.getPlaybookId(), retCode, m.Meta.Trace)

	M := iMessage()
	M.Meta = meta{
		// replace type, just record this message
		MessageType: ToRet,
		Trace:       bury.Meta.getTrace(),
		// replace seq, if we need JUMP node seq inc one after entering this, using bury,
		// on the contrary, using m instead.
		Sequence:  bury.Meta.getSequence(),
		Topic:     bury.Meta.getTopic(),
		Signature: bury.Meta.getSignature(),
		Timestamp: cast.ToString(time.Now().UnixNano() / 1e6),
		// do not pop stack again!
		Domain:      m.Context.getStack().getCurrDomain(),
		Mod:         bury.Meta.getMod(),
		ServiceAddr: bury.Meta.getServiceAddr(),
		Sync:        bury.Meta.Sync,
	}

	M.Task = task{
		// kidnap and replace id by the original one
		PlaybookId: pb.GetId(),
		SnapshotId: pb.GetSnapshotId(),
		// avoiding infinite loops
		NodeCode: retCode,
		Status:   bury.Task.getStatus(),
		Input:    last,
		// define internal sealed playbook's output.
		Output:        m.Task.getOutput(),
		Xrn:           bury.Task.getXrn(),
		Configuration: bury.Task.GetConfiguration(),
	}

	M.Context = bury.Context
	return M, nil
}

// ToSlave or ToRet
func (m *message) handleTrace(output interface{}) error {
	if pb, ok := gCache.GetPlaybook(m.Task.PlaybookId); ok {

		// resolve versions can visit node which is not exist!
		n, ok := pb.GetNode(m.Task.NodeCode)
		if !ok {
			return fmt.Errorf("node %v not exist, slot %v", m.Task.NodeCode, m.Meta.Trace)
		}

		if !m.isAbort() {
			// there are concurrent situation， cause weak and mark using different lock.
			if n.isWeak() {
				//logx.L().Debugf("handleTrace cache weak slot:%v, node:%v", m.Meta.Trace, m.Task.NodeCode)
				m.cacheTrace(nil, Weak)
			}

			// in the concurrent case, when message abort we should not cacheTrace anymore
			if n.isMark() {
				//logx.L().Debugf("handleTrace cache mark, slot:%v, node:%v", m.Meta.Trace, m.Task.NodeCode)
				m.cacheTrace(output, Mark)
				m.cacheTrace(nil, Ctx)
			}
		}
		return nil
	}

	return fmt.Errorf("palybook %v not exist, slot %v", m.Task.PlaybookId, m.Meta.Trace)
}

// NewServerMessage the worker response
// to support cases where a node must be completed simultaneously
// by multiple parents. the inlet gives another solution to solving ack states.
// m already executed. so the state of the the node  should be marked in the
// redis-cache for the sake of slot the ack result.
// so we can support the node expecting multiple parent node under executed.
// redis k-v should be:
// slot id: [node1, node2, node3] which slice is pre-filled by
// the validate period. the key should be revoke while pb stop.
// execute handler this
func (m *message) NewServerMessage(output interface{}, env Env, hb *HeartBeat) *message {
	//logx.L().Debugf("ToServer, slot:%v, node:%v", m.Meta.Trace, m.Task.getNodeCode())
	M := iMessage()
	// add sequence
	m.Meta.addSequence()
	// set hb
	m.Context.setHeartBt(hb)
	// record chain
	m.Context.recordChain(m.Task.getNodeCode())
	m.Context.Env = env

	M.Meta = meta{
		MessageType: ToServer,
		Trace:       m.Meta.getTrace(),
		// no need recording this, server convert only.
		Sequence:    m.Meta.getSequence(),
		Topic:       m.Meta.getTopic(),
		Signature:   m.Meta.getSignature(),
		Timestamp:   cast.ToString(time.Now().UnixNano() / 1e6),
		Domain:      m.Meta.getDomain(),
		Mod:         m.Meta.getMod(),
		ServiceAddr: m.Meta.getServiceAddr(),
		Sync:        m.Meta.Sync,
	}

	M.Task = task{
		PlaybookId:    m.Task.getPlaybookId(),
		SnapshotId:    m.Task.getSnapshotId(),
		NodeCode:      m.Task.getNodeCode(),
		Status:        Executed,
		Input:         m.Task.getInput(),
		Output:        output,
		Xrn:           m.Task.getXrn(),
		Configuration: m.Task.Configuration,
	}

	M.Context = m.Context
	return M
}

// NewExceptionMessage is the interface for worker response err.
func (m *message) NewExceptionMessage(output interface{}) *message {
	metric.Metric(tag.WorkerException)
	m.Meta.addSequence()
	m.Context.recordChain(m.Task.getNodeCode())
	var err = fmt.Errorf("NewExceptionMessage code:%v, task:%v, err:%v", m.Task.getNodeCode(), m.Task.Xrn.TaskInfo, output)
	M := iMessage()
	defer M.deferErrorHandler(err)

	M.Meta = meta{
		MessageType: ToException,
		Trace:       m.Meta.getTrace(),
		Sequence:    m.Meta.getSequence(),
		Topic:       m.Meta.getTopic(),
		Signature:   m.Meta.getSignature(),
		Timestamp:   cast.ToString(time.Now().UnixNano() / 1e6),
		Domain:      m.Meta.getDomain(),
		Mod:         m.Meta.getMod(),
		ServiceAddr: m.Meta.getServiceAddr(),
		Sync:        m.Meta.Sync,
	}

	M.Context = m.Context
	M.Context.setException(exception{Content: output})

	M.Task = task{
		PlaybookId: m.Task.getPlaybookId(),
		SnapshotId: m.Task.getSnapshotId(),
		NodeCode:   m.Task.getNodeCode(),
		Status:     Fail,
		// when exception happens, input should be empty? output?
		Input:         m.Task.getInput(),
		Output:        output,
		Xrn:           nil,
		Configuration: m.Task.GetConfiguration(),
	}
	return M
}

// (NewExceptionMessage, setExceptionErrorForNode, setExceptionError) will call to yield sync and
// abort process task.
func (m *message) deferErrorHandler(err error) {
	defer func() {
		eve := Event{
			topic:      m.Meta.Topic,
			TraceId:    m.Meta.Trace,
			PlaybookId: m.Task.PlaybookId,
			Body:       err.Error(),
			Time:       time.Now(),
		}
		eve.Dispatch()
	}()

	logx.L().Infof("process aborts, slot %v err:%v", m.Meta.Trace, err)
	m.abort()

	if m.Meta.needSync() {
		//KV.RedisKv.KVSet(cast.ToString(m.Meta.Trace), 1, 10*time.Minute)
		m.yieldSync(func() (interface{}, error) {
			return nil, err
		})
	}

	if pb, ok := gCache.GetPlaybook(m.Task.getPlaybookId()); ok && pb.hasMark() {
		m.delTrace()
	}
}

// setExceptionError message internal err
func (m *message) setExceptionError(err error) {
	defer m.deferErrorHandler(fmt.Errorf("setExceptionError code:%v, task:%v, err:%v", m.Task.getNodeCode(), m.Task.Xrn.TaskInfo, err))
	metric.Metric(tag.ServerException)
	m.Meta.setMessageType(ToException)
	m.Task.setStatus(Fail)
	m.Context.setException(exception{Content: err.Error()})
}

// setExceptionError message internal err
func (m *message) setExceptionErrorForNode(nodeCode string, err error) *message {
	metric.Metric(tag.ServerException)
	m.Meta.addSequence()
	m.Context.recordChain(nodeCode)

	M := iMessage()
	defer M.deferErrorHandler(fmt.Errorf("setExceptionErrorForNode code:%v, task:%v, err:%v", m.Task.getNodeCode(), m.Task.Xrn.TaskInfo, err))

	M.Meta = meta{
		MessageType: ToException,
		Trace:       m.Meta.getTrace(),
		Sequence:    m.Meta.getSequence(),
		Topic:       m.Meta.getTopic(),
		Signature:   m.Meta.getSignature(),
		Timestamp:   cast.ToString(time.Now().UnixNano() / 1e6),
		Domain:      m.Meta.getDomain(),
		Mod:         m.Meta.getMod(),
		ServiceAddr: m.Meta.getServiceAddr(),
		Sync:        m.Meta.Sync,
	}

	M.Context = m.Context
	M.Context.setException(exception{Content: err.Error()})

	M.Task = task{
		PlaybookId: m.Task.getPlaybookId(),
		SnapshotId: m.Task.getSnapshotId(),
		NodeCode:   nodeCode,
		Status:     Fail,
		// when exception happens, input should be empty? output?
		Input:         m.Task.getInput(),
		Output:        err.Error(),
		Xrn:           nil,
		Configuration: nil,
	}
	return M
}

// The order in which messages are sent is no longer just fifO,
// but depends on the heap to achieve priority,
// synchronous messages will affect faster
func (m *message) Dispatch() {
	dispatcher.lock.Lock()
	heap.Push(dispatcher.msgQ, m)
	dispatcher.lock.Unlock()
	semaphore <- struct{}{}
}

// Consider ambiguous when a flying message encounter a updated snapshot.
// each message handled by ToServer or ToException should check snaps by
// ResolveAmbiguous before hatch new messages.
// this options turn the cluster to the state of final consistency.
func (m *message) ResolveAmbiguous() error {
	pbId := m.Task.getPlaybookId()
	lastSnapId := m.Task.getSnapshotId()
	pb, ok := gCache.GetPlaybook(pbId)
	if !ok {
		return fmt.Errorf("invalid pb id")
	}

	currSnapId := pb.GetSnapshotId()
	// when shit happens
	if lastSnapId != currSnapId {
		res, err := findPb(pbId)
		if err != nil {
			return err
		}
		if res.SnapshotId == currSnapId {
			logx.L().Infof("version control for pb: [%d] failed! receive snapshot: [%d], current snapshot: [%d], msgtype:%v",
				pbId, lastSnapId, currSnapId, m.Meta.MessageType.String())
			// indicates last snap laps behind real snap
			emsg := m.setExceptionErrorForNode(m.Task.NodeCode,
				fmt.Errorf("version control for pb: [%d] failed! receive snapshot: [%d], current snapshot: [%d]",
					pbId, lastSnapId, currSnapId))
			emsg.Dispatch()
			return fmt.Errorf("version control update fail")
		} else if res.SnapshotId == lastSnapId {
			// if current snap lapse, reload please
			logx.L().Infof("ReloadSnapshot now, lastSnapId:%v, pbId:%v", lastSnapId, pbId)
			err := ReloadSnapshot(lastSnapId, pbId)
			if err != nil {
				logx.L().Errorf("ReloadSnapshot err %v", err)
				return err
			}
		} else {
			// oops, mess up
			logx.L().Errorf("ResolveAmbiguous", "messed up")
			err = fmt.Errorf("slot messed up :(")
			return err
		}
	}

	return nil
}

func (m *message) WaitSync(c ctx.Context) ([]byte, error) {
	if m.Meta.needSync() {
		// define expire
		cc, _ := ctx.WithTimeout(c, 60*time.Second)
		c := make(waitSync)
		// cache sync & delete it while code return
		m.cacheSync(c)
		defer m.deleteSync()

		select {
		case f := <-c:
			res, err := f()
			if err != nil {
				logx.L().Infof("WaitSync receive err, slot：%v", m.Meta.Trace)
				return nil, err
			}
			b, err := json.Marshal(res)
			return b, err

		case <-cc.Done():
			// abort here but must be clearly recognize that
			// the flow may be execute at the same time.
			// so there are some concurrency data matter should be
			// noticed.
			m.abort()

			if pb, ok := gCache.GetPlaybook(m.Task.getPlaybookId()); ok && pb.hasMark() {
				logx.L().Debugf("ctx done, do delTrace")
				m.delTrace()
			}
			//logx.L(zap.Int(logx.PlaybookId, m.Task.PlaybookId)).Debugf("WaitSync timeout for slot:%v", m.Meta.Trace)
			// we should metric exception here.
			metric.Metric(tag.ServerException)
			return nil, fmt.Errorf("sync result message wait 10 sec exceed, but no response yet")
		}
	}

	return nil, nil
}

type Event struct {
	topic      string    `json:"-"`
	TraceId    int64     `json:"trace_id"`
	PlaybookId int       `json:"playbook_id"`
	Body       string    `json:"body"`
	Time       time.Time `json:"time"`
}

func (e *Event) GetTopic() string {
	return e.topic
}

// Dispatch event handler
func (e *Event) Dispatch() {
	dispatcher.schemeC <- e
}

func (e *Event) toReader() io.Reader {
	return bytes.NewBuffer(e.toBytes())
}

func (e *Event) toBytes() []byte {
	body, _ := json.Marshal(e)
	body = append(body, EVENT)
	return body
}

type Msg struct {
	// message getCacheKey result
	Key         AcKey      `json:"key"`
	ServiceAddr string     `json:"srvAddr"`
	Signature   *Signature `json:"Sign"`

	Input         interface{}    `json:"input"`
	Configuration *configuration `json:"config"`
	Env           Env            `json:"env"`
	Time          time.Time      `json:"time"`
}

// Dispatch event handler
func (m *Msg) Dispatch() {
	dispatcher.schemeC <- m
}

func (m *Msg) toReader() io.Reader {
	return bytes.NewBuffer(m.toBytes())
}

func (m *Msg) toBytes() []byte {
	body, _ := json.Marshal(m)
	body = append(body, MESSAGE)
	return body
}

type signType int

const (
	Ecc signType = iota
	Jwt
)

// Tamper proofing Sign for message.
func (m *Msg) Sign(st signType) {
	switch st {
	case Jwt:
		var (
			jw  string
			err error
		)

		// Sign the hash using jw
		jw, err = jwtutil.GenTokenString(m.tMd5())
		if err != nil {
			logx.L().Errorf("Sign error %v", err)
			jw = ""
		}

		m.Signature = &Signature{
			Type: st,
			Jwt:  jw,
			Hash: toolutil.String2Byte(toolutil.CheckSum(toolutil.String2Byte(jw))),
		}
	case Ecc:
		pt := m.getPlain()
		sig, _ := eccutil.EccSign(pt)
		m.Signature = &Signature{
			Type: st,
			Jwt:  "",
			Hash: sig,
		}
	default:
		panic("not implement yet")
	}

}

func (m *Msg) CheckSign() error {
	typ, jt, md5x := m.Signature.Type, m.Signature.Jwt, m.Signature.Hash
	switch typ {
	case Jwt:
		if b := toolutil.CheckSum(toolutil.String2Byte(jt)); b != toolutil.Bytes2string(md5x) {
			return fmt.Errorf("validate signature hash failed")
		}
		token, err := jwtutil.CheckTokenString(jt)
		if err != nil {
			return fmt.Errorf("validate signature jwt failed")
		}

		tmd5jwt := token.Claims.(jwt.MapClaims)
		if md5, ok := tmd5jwt["md5"]; ok && m.tMd5()["md5"] == md5 {
			return nil
		}
		return fmt.Errorf("manipulated message checked by jwt md5, real: %v, curr:%v", tmd5jwt["md5"], m.tMd5()["md5"])
	case Ecc:
		pt := m.getPlain()
		if !eccutil.EccSignVer(pt, m.Signature.Hash) {
			return fmt.Errorf("incorrect sign")
		}

	default:
		return fmt.Errorf("invalide sign type, %v", typ)
	}

	return nil
}

func (m *Msg) tMd5() map[string]interface{} {
	// message md5
	md5 := toolutil.CheckSum(m.getPlain())
	return map[string]interface{}{"md5": md5}
}

func (m *Msg) getPlain() []byte {
	return toolutil.String2Byte(string(m.Key) + m.ServiceAddr)
}
