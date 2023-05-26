// Author: huaxr
// Time:   2021/8/12 上午11:57
// Git:    huaxr

package core

import (
	"encoding/json"
	"fmt"

	"github.com/huaxr/magicflow/component/express/parser"
	"github.com/huaxr/magicflow/component/logx"
)

type (
	nodeType string
	keyType  string
	bitMask  uint8
)

const (
	// Mark represents this node has in-degree direction to a multi input node.
	// when is true, this node output & slot id should be marked.
	//       A       B
	//       |       |        => A & B mark true
	//       |-> C <-|
	// there are two points that we should consider:
	// 1. slot cache
	// 2. merge Ctx
	markState = bitMask(1)

	// this property marks that the node does not need to be set in dependencies field.
	weakState = bitMask(2)

	// label end means no outlet for this node. we need record it to cache with output.
	endState = bitMask(4)
)

const (
	CommonNode nodeType = "node"
	Playbook   nodeType = "playbook"

	StartNodeCode = ""
	StopNodeCode  = "?"
)

type nPrivate interface {
	setDependencies(d []string)
	setWDependencies(d []string)
	setMark()
	setWeak()

	getDependencies() []string
	getWDependencies() []string
	isMark() bool
	isWeak() bool

	isEnd() bool
	setEnd()

	isHook() bool
	setHook()
	newTtl() ttl
	getTtl() ttl
	setTtl(t ttl)
	setBranchChoice(branches)
	setBranchParallel(branches)

	getCall() int
	getRet() int
	getNamespace() string
}

type NodeImpl interface {
	nPrivate
	GetNodeCode() string
	GetBranchChoice() branches
	GetBranchParallel() branches
	getNodeType() nodeType
	GetTask() *task
}

type (
	commands []struct {
		Type       string                 `json:"type"`
		Parameters map[string]interface{} `json:"parameters"`
	}
	// configuration is the task elementary field to analyzed
	// the task should reference.
	configuration struct {
		// when worker executes fail, retry
		Retry         int      `json:"retry,omitempty"`
		Timeout       int      `json:"timeout,omitempty"`
		BeforeExecute commands `json:"before_execute,omitempty"`
		AfterExecute  commands `json:"after_execute,omitempty"`
		// when calling remote node, the token should be filled
		Token string `json:"token,omitempty"`

		// keys just denote and handleTrace.
		// {"number1":"int", "number2":"int"}
		InputKeys map[string]keyType `json:"input_keys,omitempty"`
		// {"number":"int"}
		OutputKeys map[string]keyType `json:"output_keys,omitempty"`
	}
)

//go:generate go run github.com/fjl/gencodec -type node -out node_gen.go
// node implied there is a task to be executed, whether
// proxy currently or future sustained cloud function from
// tal-architecture. we can define xrn instead of using a label
// to refer to our client method just consider the expand for the task.
type node struct {
	// Branches[0] represents choice.
	// node has at least one edge with conditions
	// the branches must execute parallel node may has more
	// than one branches to be handled by parallel

	// Branches[1] represents parallel.
	// constraint is differ from branches, which is a precondition involved
	// to express the execution strategy of subsequent nodes. restrains itself
	// before execution check.
	/*   +---B---+
		 |       |(connection1)
	A--->|       |---> D
		 |       |(connection2)
		 +---C---+           */
	// D node should be executed after B&C done already.
	// the execution can continue only if both its parent nodes are in the
	// completed execution state.

	// e.g. B's Constraint is [D], C's Constraint is [D]
	// D's Constraint is empty, but it's Dependence field is "B,C"
	// keep abreast with the restrains node output.
	// resign the node to the previous execute result instead of respective
	// execute without restrict.

	// note that conditional branches must be mutually exclusive,
	// otherwise there will be ambiguity, it's allowed users manger
	// their playbook with non exclusive conditions, but only one will be selected
	// at execution time (Hatch)
	Branches map[branchType]branches `json:"branches" gencodec:"required"`

	Task *task `json:"task"`

	/*      A                 Let's assume that AC and AD are parallel sides
		 /  |   \             AB is the conditional expression edge
		/   =    =            the R node requires the output of B,C and D1
	   B    C     D           so we need to combine storage and KV to make this progenitor relationship work
	   |    |    /   \          1. When C and D1 are executed, the resulting state is stored in KV
	   |    |   D1    D2        2. Since the execution output already exists in storage, context can be assembled with query
	    \   |  /                3. BR path check before deliver the message to the worker
	     \	| /
		 R(B & C & D1)               if exist keys, go through it, nor sending to __internal queue.
	The overall architectural design ushered in a major overhaul.
	we using GRPc dialing the mod of traceId, so a chain of execution must run on a specific pod.
	which means we can discard redis-KV and just record

	R's `Dependencies` is [B,C,D1] `Mark` is false
	B、C、D1's `Dependencies` is [] `Mark` is true

	weakDepends performance is consistent with dependencies which notes the choice edge's next code
	field, it's init at load time, and the corresponding thing is weak filed.
	*/
	dependencies, weakDepends depends

	// hook is related to node, it's not belong to task.
	// hook after task executed and to server happen.
	Hook bool `json:"hook"`

	NodeType nodeType `json:"node_type"`
	// NodeCode use node code to identify flow. this field couldn't has prefix contains $/_/./
	NodeCode string `json:"node_code"`

	// state is the combination of mark/weak/end
	state bitMask

	// tell slime that this node should not be delete.
	ttl ttl
}

func (n *node) getNodeType() nodeType {
	return n.NodeType
}

func (n *node) GetNodeCode() string {
	return n.NodeCode
}

// only consider adjacent points and edges
func (n *node) newTtl() ttl {
	var tt = make(ttl)
	tt = ttlx(n.Task.InputAttribute, tt)
	// edge's express also cite by $., should be merged here.
	n.setTtl(merge(tt, n.Branches[BranchChoice].getTtl()))
	return n.getTtl()
}

func (n *node) getTtl() ttl {
	return n.ttl
}

func (n *node) setTtl(t ttl) {
	n.ttl = t
}

type pbNode struct {
	node
	Call int `json:"call"`
	// for restoring the current(parent) message.
	Ret int `json:"ret"`
}

// calling other app's worker handler
type rpbNode struct {
	pbNode
	// calling app
	Namespace string `json:"namespace"`
}

func newLocalPbNode(code string, retPbId int, task *task) NodeImpl {
	//if task.Xrn.getCallId() == retPbId {
	//	logx.L().Errorf("newLocalPbNode", "could not set the same id")
	//	return nil
	//}
	n := new(pbNode)
	n.NodeType = Playbook
	n.NodeCode = code
	n.Branches = map[branchType]branches{}
	//n.BranchChoice = make([]*edge, 0)
	//n.BranchParallel = make([]*edge, 0)
	n.Task = task
	n.Call = task.Xrn.getCallId()
	n.Ret = retPbId
	return n
}

func newRemotePbNode(code string, retPbId int, task *task) NodeImpl {
	n := new(rpbNode)
	n.NodeType = Playbook
	n.NodeCode = code
	n.Branches = map[branchType]branches{}

	n.Task = task
	n.Call = task.Xrn.getCallId()
	n.Ret = retPbId
	n.Namespace = task.Xrn.Namespace
	return n
}

// newNode creat a normal node to handle js.
func newNode(code string, task *task) NodeImpl {
	if code == StartNodeCode || code == StopNodeCode {
		logx.L().Errorf("newNode", "the code already used!")
		return nil
	}
	n := new(node)
	n.NodeType = CommonNode
	n.NodeCode = code
	n.Branches = map[branchType]branches{}
	task.NodeCode = code
	n.Task = task
	return n
}

func newNopNode(code string) NodeImpl {
	n := new(node)
	n.NodeType = CommonNode
	n.NodeCode = code
	n.Branches = map[branchType]branches{}
	n.Task = nil
	return n
}

func (n *node) GetTask() *task {
	return n.Task
}

func (n *rpbNode) getCall() int {
	return n.Call
}

func (n *rpbNode) getRet() int {
	return n.Ret
}

func (n *rpbNode) getNamespace() string {
	return n.Namespace
}

func (n *pbNode) getCall() int {
	return n.Call
}

func (n *pbNode) getRet() int {
	return n.Ret
}

func (n *pbNode) getNamespace() string {
	panic("")
}

func (n *node) getCall() int {
	panic("")
}

func (n *node) getRet() int {
	panic("")
}

func (n *node) getNamespace() string {
	panic("")
}

func (n *node) GetBranchChoice() branches {
	return n.Branches[BranchChoice]
}

func (n *node) GetBranchParallel() branches {
	return n.Branches[BranchParallel]
}

func (n *node) setBranchChoice(e branches) {
	n.Branches[BranchChoice] = e
}

func (n *node) setBranchParallel(e branches) {
	n.Branches[BranchParallel] = e
}

func (n *node) setDependencies(d []string) {
	n.dependencies = d
}

func (n *node) setWDependencies(d []string) {
	n.weakDepends = d
}

func (n *node) getDependencies() []string {
	return n.dependencies
}

func (n *node) getWDependencies() []string {
	return n.weakDepends
}

func (n *node) isMark() bool {
	return n.state&markState == markState
}

func (n *node) isWeak() bool {
	return n.state&weakState == weakState
}

func (n *node) isEnd() bool {
	return n.state&endState == endState
}

func (n *node) setEnd() {
	//n.end = true
	n.state |= endState
}

func (n *node) setMark() {
	n.state |= markState
}

func (n *node) setWeak() {
	n.state |= weakState
}

func (n *node) isHook() bool {
	return n.Hook
}

func (n *node) setHook() {
	n.Hook = true
}

func (n *node) UnmarshalJSON(b []byte) error {
	var (
		objMap map[string]*json.RawMessage
		nt, nc string
		bs     map[branchType]branches
		t      task
		hook   bool
	)

	err := parser.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	NodeTyper, ok := objMap["node_type"]
	if !ok {
		return fmt.Errorf("missing field nodeType")
	}

	NodeCode, ok := objMap["node_code"]
	if !ok {
		return fmt.Errorf("missing field NodeCode")
	}

	Branches, ok := objMap["branches"]
	if !ok {
		return fmt.Errorf("missing field Branches")
	}

	Task, ok := objMap["task"]
	if !ok {
		return fmt.Errorf("missing field Task")
	}

	Hook, ok := objMap["hook"]
	if !ok {
		return fmt.Errorf("missing field Task")
	}

	if NodeCode == nil {
		return fmt.Errorf("NodeCode field must not be nil")
	}

	if err := parser.Unmarshal(*NodeTyper, &nt); err != nil {
		return fmt.Errorf("unmarshal nodeType err: %v", err)
	}

	if err := parser.Unmarshal(*NodeCode, &nc); err != nil {
		return fmt.Errorf("unmarshal NodeCode err: %v", err)
	}

	if err := parser.Unmarshal(*Branches, &bs); err != nil {
		return fmt.Errorf("unmarshal Branches err: %v", err)
	}

	if err := parser.Unmarshal(*Task, &t); err != nil {
		return fmt.Errorf("unmarshal Task err: %v", err)
	}

	if err := parser.Unmarshal(*Hook, &hook); err != nil {
		return fmt.Errorf("unmarshal Hook err: %v", err)
	}

	n.NodeType = nodeType(nt)
	n.NodeCode = nc
	n.Branches = bs
	n.Task = &t
	n.Hook = hook
	return nil
}

func (n *pbNode) UnmarshalJSON(b []byte) error {
	var (
		objMap    map[string]*json.RawMessage
		nt, nc    string
		call, ret int
		bs        map[branchType]branches
		t         task
		hook      bool
	)

	err := parser.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	NodeTyper, ok := objMap["node_type"]
	if !ok {
		return fmt.Errorf("missing field nodeType")
	}

	NodeCode, ok := objMap["node_code"]
	if !ok {
		return fmt.Errorf("missing field NodeCode")
	}

	Branches, ok := objMap["branches"]
	if !ok {
		return fmt.Errorf("missing field Branches")
	}

	Task, ok := objMap["task"]
	if !ok {
		return fmt.Errorf("missing field Task")
	}

	Hook, ok := objMap["hook"]
	if !ok {
		return fmt.Errorf("missing field Hook")
	}

	Call, ok := objMap["call"]
	if !ok {
		return fmt.Errorf("missing field Call")
	}

	Ret, ok := objMap["ret"]
	if !ok {
		return fmt.Errorf("missing field Ret")
	}

	if NodeCode == nil {
		return fmt.Errorf("NodeCode field must not be nil")
	}

	if err := parser.Unmarshal(*NodeTyper, &nt); err != nil {
		return fmt.Errorf("unmarshal nodeType err: %v", err)
	}

	if err := parser.Unmarshal(*NodeCode, &nc); err != nil {
		return fmt.Errorf("unmarshal NodeCode err: %v", err)
	}

	if err := parser.Unmarshal(*Branches, &bs); err != nil {
		return fmt.Errorf("unmarshal Branches err: %v", err)
	}

	if err := parser.Unmarshal(*Task, &t); err != nil {
		return fmt.Errorf("unmarshal Task err: %v", err)
	}

	if err := parser.Unmarshal(*Call, &call); err != nil {
		return fmt.Errorf("unmarshal Call err: %v", err)
	}

	if err := parser.Unmarshal(*Ret, &ret); err != nil {
		return fmt.Errorf("unmarshal Ret err: %v", err)
	}

	if err := parser.Unmarshal(*Hook, &hook); err != nil {
		return fmt.Errorf("unmarshal Hook err: %v", err)
	}

	n.NodeType = nodeType(nt)
	n.NodeCode = nc
	n.Branches = bs
	n.Task = &t
	n.Call = call
	n.Ret = ret
	n.Hook = hook
	return nil
}

func (n *rpbNode) UnmarshalJSON(b []byte) error {
	var (
		objMap    map[string]*json.RawMessage
		nt, nc    string
		call, ret int
		np        string
		bs        map[branchType]branches
		t         task
		hook      bool
	)

	err := parser.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	NodeTyper, ok := objMap["node_type"]
	if !ok {
		return fmt.Errorf("missing field nodeType")
	}

	NodeCode, ok := objMap["node_code"]
	if !ok {
		return fmt.Errorf("missing field NodeCode")
	}

	Branches, ok := objMap["branches"]
	if !ok {
		return fmt.Errorf("missing field Branches")
	}

	Task, ok := objMap["task"]
	if !ok {
		return fmt.Errorf("missing field Task")
	}

	Hook, ok := objMap["hook"]
	if !ok {
		return fmt.Errorf("missing field Hook")
	}

	Call, ok := objMap["call"]
	if !ok {
		return fmt.Errorf("missing field Call")
	}

	Ret, ok := objMap["ret"]
	if !ok {
		return fmt.Errorf("missing field Ret")
	}

	Namespace, ok := objMap["namespace"]
	if !ok {
		return fmt.Errorf("missing field AppName")
	}

	if NodeCode == nil {
		return fmt.Errorf("NodeCode field must not be nil")
	}

	if err := parser.Unmarshal(*NodeTyper, &nt); err != nil {
		return fmt.Errorf("unmarshal nodeType err: %v", err)
	}

	if err := parser.Unmarshal(*NodeCode, &nc); err != nil {
		return fmt.Errorf("unmarshal NodeCode err: %v", err)
	}

	if err := parser.Unmarshal(*Branches, &bs); err != nil {
		return fmt.Errorf("unmarshal Branches err: %v", err)
	}

	if err := parser.Unmarshal(*Task, &t); err != nil {
		return fmt.Errorf("unmarshal Task err: %v", err)
	}

	if err := parser.Unmarshal(*Call, &call); err != nil {
		return fmt.Errorf("unmarshal Call err: %v", err)
	}

	if err := parser.Unmarshal(*Ret, &ret); err != nil {
		return fmt.Errorf("unmarshal Ret err: %v", err)
	}

	if err := parser.Unmarshal(*Namespace, &np); err != nil {
		return fmt.Errorf("unmarshal AppName err: %v", err)
	}

	if err := parser.Unmarshal(*Hook, &hook); err != nil {
		return fmt.Errorf("unmarshal Hook err: %v", err)
	}

	n.NodeType = nodeType(nt)
	n.NodeCode = nc
	n.Branches = bs
	n.Task = &t
	n.Call = call
	n.Ret = ret
	n.Namespace = np
	n.Hook = hook
	return nil
}
