// Author: huaxr
// Time:   2021/6/30 下午2:34
// Git:    huaxr

package core

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/pkg/toolutil"
)

type (
	// illuminate: if stack using message the error when happen like below.
	// json: unsupported value: encountered a cycle via *global.message
	// cause golang dose not support circle struct marshal.
	// the harmony achievement of context breakpoint.
	// [2][]byte elem1 is message body, elem2 is the ret code relatively.
	stack [][2][]byte
)

func (s stack) size() int {
	return len(s)
}

func (s *stack) push(m [2][]byte) {
	*s = append(*s, m)
}

// Usage situation:
// 1: when need ret from the pb in message.ret
// 2: when exception happens in the pb in record
func (s *stack) pop() (*message, string) {
	if s.size() == 0 {
		return nil, ""
	}
	m := (*s)[len(*s)-1]
	var msg message
	err := json.Unmarshal(m[0], &msg)
	if err != nil {
		logx.L().Errorf("stack.pop", "Unmarshal message error")
		return nil, ""
	}
	*s = (*s)[:len(*s)-1]
	return &msg, toolutil.Bytes2string(m[1])
}

// Stack getCurrDomain return the slice second param which define
// the pre push content.
func (s stack) getCurrDomain() string {
	if len(s) == 0 {
		return ""
	}
	last := (s)[len(s)-1]
	return toolutil.Bytes2string(last[1])
}

type Env map[string]interface{}

func (c Env) Get(key string) (val interface{}, exist bool) {
	val, exist = c[key]
	return
}

func (c Env) Set(key string, val interface{}) {
	c[key] = val
}

// Deadline always returns that there is no deadline (ok==false),
// maybe you want to use Request.Context().Deadline() instead.
func (c *Env) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done always returns nil (chan which will wait forever),
// if you want to aborts your work when the connection was closed
// you should use Request.Context().Done() instead.
func (c *Env) Done() <-chan struct{} {
	return nil
}

// Err always returns nil, maybe you want to use Request.Context().Err() instead.
func (c *Env) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Env) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}

type HeartBeat struct {
	Alive    string `json:"alive"`
	CpuUsage string `json:"cpu_usage"`
	MemIdle  string `json:"mem_idle"`
	MemAlloc string `json:"mem_alloc"`
	Host     string `json:"host"`
	OS       string `json:"os"`
}

type exception struct {
	Content interface{} `json:"content"`
}

func (e exception) GetContent() interface{} {
	return e.Content
}

// Context is a interim storage for a successive flow to heritage
// from previous message context, in which we can store information,
// and increase the interaction abilities.
type context struct {
	// desperate retrospective pointer to the last nodeName.
	// (albeit reluctantly, it retrench the message body.)
	// namely substitute resent body of slave output.
	Last interface{} `json:"last"`
	// Trigger:body -> AAA:output, BBB:output
	// vague trigger key.
	Store     map[string]interface{} `json:"store"`
	Exception exception              `json:"exception"`
	// for restoring the parent message context like esp and ebp register dose.
	// parcel with the message (maybe omit) and give a shit.
	// [[2], [2]...] when exception happens, a series of ret to record the chain.
	Stack stack `json:"stack,omitempty"`
	Env   Env   `json:"Ctx"`
	// A->B->C... represents flow direction, avoiding multi inlet to one
	// node
	Chain []string `json:"chain,omitempty"`

	Heartbeat *HeartBeat `json:"heartbeat"`
}

func (c *context) Get(key string) (val interface{}, exist bool) {
	return c.Env.Get(key)
}

func (c *context) Set(key string, val interface{}) {
	c.Env.Set(key, val)
}

func (c *context) GetEnv() Env {
	return c.Env
}

// getException
func (c *context) getException() exception {
	return c.Exception
}

func (c *context) setException(e exception) {
	c.Exception = e
}

func (c *context) setHeartBt(hb *HeartBeat) {
	c.Heartbeat = hb
}

func (c *context) getStack() *stack {
	return &c.Stack
}

func (c *context) setStore(s map[string]interface{}) {
	c.Store = s
}

func (c *context) setEnv(s map[string]interface{}) {
	c.Env = s
}

// ret, toServer, toExcept, double-record
func (c *context) recordChain(nodeCode string) {
	c.Chain = append(c.Chain, nodeCode)
}

func (c *context) getChain() string {
	return strings.Join(c.Chain, "->")
}

func (c *context) setLast(s interface{}) {
	c.Last = s
}

// setSlaveOutput fill the context with node key-output pair, the trigger body
// will padding the context filed by the _internal_trigger key.
func (c *context) setSlaveOutput(nodeCode string, output interface{}) {
	// bypass saving START node of trigger body, because the NewTriggerMessage
	// has already keep it.
	if nodeCode == StartNodeCode || nodeCode == StopNodeCode {
		return
	}

	c.Last = output
	if c.Store == nil {
		c.Store = make(map[string]interface{})
	}
	c.Store[nodeCode] = output
}

func (c *context) getStore() map[string]interface{} {
	return c.Store
}

func (c *context) getEnv() map[string]interface{} {
	return c.Env
}

func (c *context) getLast() interface{} {
	return c.Last
}
