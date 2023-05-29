// Author: huaxr
// Time:   2021/8/12 上午11:58
// Git:    huaxr

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/dispatch"
	"github.com/huaxr/magicflow/component/express/parser"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/plugin/limiter"
	"github.com/huaxr/magicflow/component/plugin/selector"
	"github.com/huaxr/magicflow/pkg/accutil"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/spf13/cast"
)

// playbook is StateMachine headquarters to a certain extent.
// you can use calling addEdge/addNode to assemble integrate structure with
// branches and nodes, once pb generated, it will register in the local gCache
// until create/update event captured in the etcd context.
// project survey view observer provided with connections illustration.
type playbook struct {
	sync.Mutex `json:"-"`
	internal   struct {
		// slime(ttl implement) calculate the worker node's output which stored in the
		// message's context temporally, if we keeping record all the slave output instead of
		// withdraw some never used node key in context, it might broaden expansion rapidly
		// so slime will record all the node's output need be store count until the flow
		// aborts. transient node output should be revoke.
		slime slimer
		// Redundancy record for incline the calling times of search first node.
		app, name, user string

		startAt string
		hasMark bool
		hasHook bool

		// app token for authorisation when connecting
		// pb token for remote calling validate
		appToken, pbToken string
	} `json:"-"`

	// index id
	Id int `json:"id"`

	// snapshot id keeps the flow dominated by the same draft simultaneously.
	SnapshotId int                 `json:"snapshot_id"`
	States     map[string]NodeImpl `json:"states"`
}

type RawPlaybook struct {
	PlaybookId int `json:"playbook_id"`
	Nodes      []struct {
		// NodeName       string      `json:"node_name"`
		NodeCode       string      `json:"node_code"`
		InputAttribute interface{} `json:"input_attribute"` // map or string
		TaskId         int         `json:"task_id"`
		Hook           int         `json:"hook"`
	} `json:"nodes"`

	Edges []struct {
		FromNode  string `json:"from_node"`
		ToNode    string `json:"to_node"`
		Condition string `json:"condition"`
		Priority  int    `json:"priority"`
	} `json:"edges"`
}

func NewPlayBook(in io.Reader) (*playbook, error) {
	var req RawPlaybook
	err := json.NewDecoder(in).Decode(&req)
	if err != nil {
		return nil, err
	}

	playbook := new(playbook)
	playbook.States = make(map[string]NodeImpl)
	playbook.Id = req.PlaybookId

	var ns = make([]models.Node, 0)
	var eg = make([]models.Edge, 0)

	for _, node := range req.Nodes {

		if len(node.NodeCode) >= 8 || len(node.NodeCode) <= 3 {
			return nil, errors.New(fmt.Sprintf("node_code validate fail: too long or too short [3:8]"))
		}

		if err := toolutil.AZaz09_(node.NodeCode); err != nil {
			return nil, fmt.Errorf("node_code validate fail:%v", err.Error())
		}

		var attr string
		switch reflect.TypeOf(node.InputAttribute).Kind() {
		case reflect.Map:
			b, _ := json.Marshal(node.InputAttribute)
			attr = toolutil.Bytes2string(b)
		case reflect.String:
			attr = node.InputAttribute.(string)
		}
		ns = append(ns, models.Node{
			NodeCode:       node.NodeCode,
			InputAttribute: attr,
			PlaybookId:     req.PlaybookId,
			TaskId:         node.TaskId,
			Hook:           node.Hook,
		})
	}

	for _, edge := range req.Edges {
		eg = append(eg, models.Edge{
			PlaybookId: req.PlaybookId,
			FromNode:   edge.FromNode,
			ToNode:     edge.ToNode,
			Condition:  edge.Condition,
			Priority:   edge.Priority,
		})
	}
	err = playbook.loadsN(ns)
	if err != nil {
		return nil, err
	}
	err = playbook.loadsE(eg)
	if err != nil {
		return nil, err
	}

	// here we calculate slime and filling the playbook
	err = playbook.Validate()
	if err != nil {
		logx.L().Errorf("NewPlayBook validate pb error, %v", err)
		return nil, err
	}
	gCache.SetPlaybook(req.PlaybookId, playbook)
	return playbook, nil
}

func SetNamespace(app *models.App) {
	var (
		eps int32 = 100
		srv       = make(map[string]string)
	)
	if int32(app.Eps) > 0 {
		eps = int32(app.Eps)
	}

	for _, s := range strings.Split(app.Brokers, ",") {
		// broker  10.187.114.211:4150 lost, there got a space leading the ip address,
		// which cause broker not found and message could not send any more.
		s = strings.TrimSpace(s)
		if dispatch.MQ(app.BrokerType) == dispatch.Kafka {
			if len(strings.Split(s, "-")) != 2 {
				logx.L().Errorf("SetNamespace kafka broker must ip-partition")
				continue
			}
		}
		srv[s] = ""
	}
	logx.L().Infof("add srv for app %v %v", app.AppName, srv)

	gCache.SetNamespace(app.AppName, &Namespace{
		limiter:  limiter.NewRateLimiter(app.AppName, eps),
		selector: selector.NewSelector(selector.RoundRobin, srv),
		enable:   1,
		share:    cast.ToBool(app.Share),
	})
}

// getStartNode
func (pb *playbook) getStartNode() NodeImpl {
	return pb.States[pb.internal.startAt]
}

func (pb *playbook) setHasMark() {
	pb.internal.hasMark = true
}

func (pb *playbook) setHasHook() {
	pb.internal.hasHook = true
}

func (pb *playbook) hasMark() bool {
	return pb.internal.hasMark
}

func (pb *playbook) HasHook() bool {
	return pb.internal.hasHook
}

func (pb *playbook) GetNodes() map[string]NodeImpl {
	return pb.States
}

func (pb *playbook) SetSnapshotId(id int) {
	pb.SnapshotId = id
}

// update all node task ssid
func (pb *playbook) updateAllNodeTaskSnapId() {
	for _, i := range pb.States {
		i.GetTask().SnapshotId = pb.SnapshotId
	}
}

func (pb *playbook) GetSnapshotId() int {
	return pb.SnapshotId
}

// GetId
func (pb *playbook) GetId() int {
	return pb.Id
}

// GetApp
func (pb *playbook) GetApp() string {
	return pb.internal.app
}

// getSlim
func (pb *playbook) getSlim() slimer {
	return pb.internal.slime
}

// GetAppToken
func (pb *playbook) GetAppToken() string {
	return pb.internal.appToken
}

// GetPbToken
func (pb *playbook) GetPbToken() string {
	return pb.internal.pbToken
}

// getNode
func (pb *playbook) GetNode(nodeName string) (n NodeImpl, ok bool) {
	n, ok = pb.States[nodeName]
	return n, ok
}

func (pb *playbook) HasNode(nodeName string) bool {
	_, ok := pb.States[nodeName]
	return ok
}

// addNode
func (pb *playbook) addNode(n NodeImpl) {
	if n == nil {
		return
	}
	pb.Lock()
	defer pb.Unlock()
	pb.States[n.GetNodeCode()] = n
}

// revokeNode remove a node form the playbook
func (pb *playbook) revokeNode(name string) {
	pb.Lock()
	defer pb.Unlock()
	delete(pb.States, name)
}

// addNodes
func (pb *playbook) addNodes(nodes ...NodeImpl) {
	for _, i := range nodes {
		if i == nil {
			return
		}
		i.GetTask().PlaybookId = pb.Id
		pb.Lock()
		pb.States[i.GetNodeCode()] = i
		pb.Unlock()
	}
}

// fillAttrs add attributes to the pb.
func (pb *playbook) fillAttrs() error {
	plb, err := findPb(pb.Id)
	if err != nil {
		logx.L().Errorf("pb err:%v", err)
		return err
	}

	app, err := findApp(plb.AppId)
	if err != nil {
		logx.L().Errorf("app err:%v", err)
		return err
	}

	pb.internal.name = plb.Name
	pb.internal.pbToken = plb.Token
	pb.internal.app = app.AppName

	pb.internal.appToken = app.Token
	pb.internal.user = app.User

	return nil

}

// loadsE find the nodeId and fill the inputAttr.
// load node information from database. we can define a template for nodes.
func (pb *playbook) loadsE(ns []models.Edge) error {
	for _, e := range ns {
		err := pb.addEdge(e.FromNode, e.ToNode, e.Condition, e.Priority)
		if err != nil {
			logx.L().Errorf("playbook.loadsE edge:%+v err:%v", e, err.Error())
			return err
		}
	}
	return nil
}

// add node to a playbook, need the node index id, and fill the
// playbook's with the record from table node and attribute.
func (pb *playbook) loadsN(ns []models.Node) error {
	for _, n := range ns {
		t, err := findTask(n.TaskId)
		if err != nil {
			logx.L().Errorf("task err:%v", err)
			return err
		}
		attr, err := newAttribute(n.InputAttribute)
		if err != nil {
			logx.L().Errorf("newAttribute err: %v", err)
			return err
		}

		// step 1: register your app and config the task sets
		// step 2: draw your playbook and binging task to one node.
		// the resource name done in this period, so the node got xesrn!
		rn := newMagicRn(t.Xrn)
		if rn == nil {
			return fmt.Errorf("bad xrn for node:%v", n.NodeCode)
		}

		// local && remote task do not need configuration field!!!
		// only normal task need.
		var config configuration
		if len(t.Configuration) > 0 {
			err = json.Unmarshal(toolutil.String2Byte(t.Configuration), &config)
			if err != nil {
				return fmt.Errorf("configuration must be a dict/map for node:%v", n.NodeCode)
			}
		}

		//
		task := &task{
			PlaybookId:     pb.Id,
			NodeCode:       n.NodeCode,
			Xrn:            rn,
			InputAttribute: attr,
			Configuration:  &config,
			Status:         Ready,
		}

		var node NodeImpl
		switch rn.getTaskType() {
		case NormalTask:
			node = newNode(n.NodeCode, task)
		case NopTask:
			node = newNopNode(n.NodeCode)
		case LocalTaskSet:
			node = newLocalPbNode(n.NodeCode, pb.Id, task)
		case RemoteTaskSet:
			node = newRemotePbNode(n.NodeCode, pb.Id, task)
		default:
			return fmt.Errorf("invalid task type for node:%v", n.NodeCode)
		}
		if node == nil {
			return errors.New("node is nil")
		}
		// set hook when it hooked
		if n.Hook == 1 {
			node.setHook()
		}

		pb.addNode(node)
	}

	return nil
}

// addEdge add a edge between two nodes which is a flow with direction.
// we call it DAG.
func (pb *playbook) addEdge(nodeCode string, nextCode string, condition string, priority int) error {
	n, ok := pb.States[nodeCode]
	if !ok {
		return fmt.Errorf("nodeCode: %v not exist, recheck please", nodeCode)
	}
	next, ok := pb.States[nextCode]
	if !ok {
		return fmt.Errorf("nextCode: %v not exist, recheck please", nextCode)
	}
	e := new(edge)
	e.NextNode = next.GetNodeCode()
	e.Express = condition
	e.Priority = priority

	if e.Express != "" {
		n.setBranchChoice(append(n.GetBranchChoice(), e))
		sort.Sort(n.GetBranchChoice())
	} else {
		n.setBranchParallel(append(n.GetBranchParallel(), e))
		sort.Sort(n.GetBranchParallel())
	}

	logx.L().Debugf("add edge nodeCode:%v, nextCode:%v success", nodeCode, nextCode)
	return nil
}

// revokeEdge
func (pb *playbook) revokeEdge(nodeName, nextName string) {
	pb.Lock()
	defer pb.Unlock()
	n, ok := pb.States[nodeName]
	if !ok {
		logx.L().Errorf("playbook.revokeEdge", "nodeId not exist")
		return
	}
	for index, i := range n.GetBranchChoice() {
		if i.NextNode == nextName {
			// bugs report
			n.setBranchChoice(append(n.GetBranchChoice()[:index], n.GetBranchChoice()[index+1:]...))
			return
		}
	}
}

// setStartNode
func (pb *playbook) setStartNode() error {
	if len(pb.States) == 1 {
		for _, i := range pb.States {
			pb.internal.startAt = i.GetNodeCode()
			break
		}
		return nil
	}

	if pb.internal.startAt != "" {
		return nil
	}

	// A->B  means B hasInDegree
	hasInDegree := make([]string, 0)
	for _, i := range pb.States {
		for _, e := range i.GetBranchChoice() {
			if !accutil.ContainsStr(hasInDegree, e.NextNode) {
				hasInDegree = append(hasInDegree, e.NextNode)
			}
		}

		for _, e := range i.GetBranchParallel() {
			if !accutil.ContainsStr(hasInDegree, e.NextNode) {
				hasInDegree = append(hasInDegree, e.NextNode)
			}
		}
	}
	noInDegree := make([]string, 0)
	for name, _ := range pb.States {
		if !accutil.ContainsStr(hasInDegree, name) {
			noInDegree = append(noInDegree, name)
		}
	}
	if len(noInDegree) != 1 {
		return fmt.Errorf("noInDegree size must be 1, current is %v", noInDegree)
	}
	pb.internal.startAt = noInDegree[0]
	return nil
}

func (pb *playbook) setTtlAndCtxParseCheck(ttyp ttlTyp) error {
	var tt slimer
	switch ttyp {
	case remember:
		tt = slimer{ttltyp: ttyp}
	case quote:
		var tts = make([]ttl, 0)
		for _, n := range pb.States {
			tx := n.newTtl()
			tts = append(tts, tx)
		}
		tt = slimer{ttltyp: ttyp, ttl: merge(tts...)}
	case dp:
		// no global ttl
		err := pb.newTtl()
		if err != nil {
			return err
		}
		tt = slimer{ttltyp: ttyp}
	}
	pb.internal.slime = tt
	return nil
}

// sets the weak reference property for each descendant node of an expression
// checks on weak nodes are avoided
func (pb *playbook) setWeakReference() error {
	for _, node := range pb.States {
		for _, edge := range node.GetBranchChoice() {
			subNodes := pb.getSubNodesOfN(pb.States[edge.NextNode], []string{})
			subNodes = append(subNodes, edge.NextNode)
			for _, nodeCode := range subNodes {
				pb.States[nodeCode].setWeak()
			}
		}
	}

	return nil
}

// find all nodes with multi input. (bigger than 1)
// then set node dependence filter and set it parent node' handleTrace true at the same time.
func (pb *playbook) setDependenceAndMark() error {
	// {"n":[n1, n2, n3], "n1":["x"]}
	var depends = map[string][]string{}
	for _, node := range pb.States {
		depends[node.GetNodeCode()] = []string{}
	}

	for _, node := range pb.States {
		for _, edge := range append(node.GetBranchParallel(), node.GetBranchChoice()...) {
			depends[edge.NextNode] = append(depends[edge.NextNode], node.GetNodeCode())
		}
	}

	for key, val := range depends {
		if len(val) > 1 {
			pb.setHasMark()

			var tmp = make([]string, 0)
			var wtmp = make([]string, 0)
			for _, i := range val {
				if !pb.States[i].isWeak() {
					tmp = append(tmp, i)
				} else {
					wtmp = append(wtmp, i)
				}
			}

			pb.States[key].setDependencies(tmp)
			pb.States[key].setWDependencies(wtmp)

			for _, v := range val {
				pb.States[v].setMark()
			}
		}
	}
	return nil
}

func (pb *playbook) setEnd() error {
	var endsAt = make([]string, 0)
	for _, node := range pb.States {
		if len(node.GetBranchChoice())+len(node.GetBranchParallel()) == 0 {
			node.setEnd()
			endsAt = append(endsAt, node.GetNodeCode())
		}
	}

	if len(endsAt) != 1 {
		return fmt.Errorf("multi end dose not alloewd")
	}
	return nil
}

func (pb *playbook) getSubNodesOfN(n NodeImpl, nodeCodes []string) []string {
	for _, e := range append(n.GetBranchChoice(), n.GetBranchParallel()...) {
		n, _ := pb.GetNode(e.NextNode)
		nodeCodes = pb.getSubNodesOfN(n, append(nodeCodes, e.NextNode))
	}

	return nodeCodes
}

// figure out all the sub nodes' ttl from this node(ancestor).
// return the ttl is not remember but a traverse result.
func (pb *playbook) traverse(n NodeImpl) ttl {
	var tt = n.getTtl()
	subNodes := pb.getSubNodesOfN(n, []string{})
	for _, code := range subNodes {
		tt = merge(tt, pb.States[code].getTtl())
	}
	return tt
}

func (pb *playbook) newTtl() error {
	// set signal node ttl
	for _, n := range pb.States {
		n.newTtl()
	}

	// traverse sub nodes and merge ttl.
	for _, n := range pb.States {
		tt := pb.traverse(n)
		n.setTtl(tt)
	}

	return nil
}

// Distinguish the DAG has endless loop. there are 3 methods normally.
// 1. Union-Find SetPlaybook
// 2. DFS
// 3. Topological sorting algorithm
func (pb *playbook) validateCircle(noInDegree []string) (err error) {
	// one node can seal a playbook.
	if len(pb.States) == 1 {
		return nil
	}
	var hasInDegree = make([]string, 0)

	// traverse states
	// get all the edges' NextNode to get a slice of hasInDegree.
	for k, i := range pb.States {
		// ignore last loop startnode, one by one check until all the node done.
		if accutil.ContainsStr(noInDegree, k) {
			continue
		}
		for _, e := range i.GetBranchChoice() {
			if !accutil.ContainsStr(hasInDegree, e.NextNode) {
				hasInDegree = append(hasInDegree, e.NextNode)
			}
		}

		for _, e := range i.GetBranchParallel() {
			if !accutil.ContainsStr(hasInDegree, e.NextNode) {
				hasInDegree = append(hasInDegree, e.NextNode)
			}
		}
	}

	// when all the nodes has income direction flows.
	// which means there is no node with no input.(the start point)
	if len(hasInDegree) == len(pb.States) || len(hasInDegree)+len(noInDegree) == len(pb.States) {
		return fmt.Errorf("acyclic validate for playbook %v err", pb.Id)
	}

	for name, _ := range pb.States {
		if accutil.ContainsStr(hasInDegree, name) || accutil.ContainsStr(noInDegree, name) {
			continue
		}
		noInDegree = append(noInDegree, name)
	}
	// gradually expand noInDegree slice until size equals states map sub 1.
	// the last noInDegree does not need loop again.

	// noInDegree is the start node currently.
	if len(noInDegree) >= len(pb.States)-1 {
		return nil
	}
	return pb.validateCircle(noInDegree)
}

// node which set hook could not used as emitted node
func (pb *playbook) validateHook() error {
	for _, node := range pb.States {
		if node.getNodeType() == Playbook {
			if node.isHook() {
				return fmt.Errorf("pb node should not set hook")
			}
			if p, ok := gCache.GetPlaybook(node.getCall()); ok && p.HasHook() {
				return fmt.Errorf("node of playbokk should not has hook")
			}
		}
	}
	return nil
}

// validateNest check if playbook node already implant into current pb.
func (pb *playbook) validateNest(pbId int) error {
	for _, i := range pb.States {
		if i.getNodeType() == Playbook {
			if i.getCall() == pbId {
				return fmt.Errorf("nest validate for playbook %v failed, because node:%v in caller pb:%v call current pb again", pbId, i.GetNodeCode(), pb.Id)
			}

			// caller id may not exist currently
			nestPb, ok := gCache.GetPlaybook(i.getCall())
			if ok {
				if err := nestPb.validateNest(pbId); err != nil {
					return err
				}
			} else {
				// maybe lapse load, ignore it
			}
		}
	}
	return nil
}

// Breadth-first traversal to validating if [mayAncestor] node is the ancestor of [curr] node.
func (pb *playbook) validateAncestor(mayAncestor NodeImpl, curr string) bool {
	for _, ancestor := range append(mayAncestor.GetBranchParallel(), mayAncestor.GetBranchChoice()...) {
		if ancestor.NextNode == curr {
			return false
		} else {
			n, _ := pb.GetNode(ancestor.NextNode)
			ok := pb.validateAncestor(n, curr)
			if !ok {
				return false
			}
		}
	}
	return true
}

// validate the hierarchy from nodes to branches.
func (pb *playbook) Validate() error {
	//circle and nest check
	err := pb.validateCircle(make([]string, 0))
	if err != nil {
		logx.L().Errorf("validateCircle err:%v", err)
		return err
	}

	// playbook not should not set hook
	err = pb.validateHook()
	if err != nil {
		logx.L().Errorf("validateHook err:%v", err)
		return err
	}

	err = pb.validateNest(pb.Id)
	if err != nil {
		logx.L().Errorf("validateNest err:%v", err)
		return err
	}

	err = pb.setStartNode()
	if err != nil {
		logx.L().Errorf("setStartNode err:%v", err)
		return err
	}

	err = pb.setTtlAndCtxParseCheck(dp)
	if err != nil {
		logx.L().Errorf("setTtlAndCtxParseCheck err:%v", err)
		return err
	}

	err = pb.setWeakReference()
	if err != nil {
		logx.L().Errorf("setWeakReference err:%v", err)
		return err
	}

	err = pb.setDependenceAndMark()
	if err != nil {
		logx.L().Errorf("setDependenceAndMark err:%v", err)
		return err
	}

	// set node attribute of end.
	err = pb.setEnd()
	if err != nil {
		logx.L().Errorf("setEnd err:%v", err)
		return err
	}

	// if node == pbNode, should check if the pb contains
	// current pb id, reduce the conflict of circle reference.
	// fill name, app, app token, limiter infos
	err = pb.fillAttrs()
	if err != nil {
		logx.L().Errorf("fillAttrs err:%v", err)
		return err
	}
	logx.L().Infof("validate success for pb:%v", pb.Id)
	return nil
}

func (pb *playbook) UnmarshalJSON(b []byte) error {
	var (
		objMap    map[string]*json.RawMessage
		id        int
		sid       int
		rawStates map[string]*json.RawMessage
	)

	err := parser.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	Id, ok := objMap["id"]
	if !ok {
		return fmt.Errorf("missing field Id")
	}

	if err := parser.Unmarshal(*Id, &id); err != nil {
		return fmt.Errorf("unmarshal id err: %v", err)
	}

	Sid, ok := objMap["snapshot_id"]
	if !ok {
		return fmt.Errorf("missing field snapshot_id")
	}
	if err := parser.Unmarshal(*Sid, &sid); err != nil {
		return fmt.Errorf("unmarshal sid err: %v", err)
	}

	States, ok := objMap["states"]
	if !ok {
		return fmt.Errorf("missing field States")
	}

	if err := parser.Unmarshal(*States, &rawStates); err != nil {
		return fmt.Errorf("unmarshal states err: %v", err)
	}

	pb.Id = id
	pb.SnapshotId = sid
	pb.States = make(map[string]NodeImpl, len(rawStates))

	for name, rawstate := range rawStates {
		if rawstate == nil {
			return fmt.Errorf("raw state for state: %s is nil", name)
		}

		var tmp map[string]*json.RawMessage
		if err := parser.Unmarshal(*rawstate, &tmp); err != nil {
			return err
		}

		var rawType string
		if err := parser.Unmarshal(*tmp["node_type"], &rawType); err != nil {
			return fmt.Errorf("unmarshal Type field for state: %s. err: %s", name, err)
		}

		var isRemote bool
		_, isRemote = tmp["namespace"]

		switch nodeType(rawType) {
		case CommonNode:
			var cn node
			if err := parser.Unmarshal(*rawstate, &cn); err != nil {
				return fmt.Errorf("unmarshal state: %s, err: %s", name, err)
			}
			pb.States[name] = &cn
		case Playbook:
			if isRemote {
				var cn rpbNode
				if err := parser.Unmarshal(*rawstate, &cn); err != nil {
					return fmt.Errorf("unmarshal state: %s, err: %s", name, err)
				}
				pb.States[name] = &cn
			} else {
				var cn pbNode
				if err := parser.Unmarshal(*rawstate, &cn); err != nil {
					return fmt.Errorf("unmarshal state: %s, err: %s", name, err)
				}
				pb.States[name] = &cn
			}
		default:
			panic("not implement yet")
		}
	}

	return nil
}

// reload snapshot from db
func reload(snapshotId, playbookId int) (*playbook, error) {
	ss, err := findSnap(snapshotId)
	if err != nil {
		return nil, err
	}
	var pb playbook
	// here snap record's snapid is 0
	err = json.Unmarshal(toolutil.String2Byte(ss.Snapshot), &pb)
	if err != nil {
		logx.L().Errorf("snapshot unmarshal is err:%+v", err)
		return nil, err
	}

	if pb.Id != playbookId {
		logx.L().Errorf("snapshot switch to the smID not same")
		return nil, errors.New("smid not same for this ssid")
	}

	err = pb.Validate()
	if err != nil {
		return nil, err
	}
	return &pb, nil
}

func ReloadSnapshot(snapshotId, playbookId int) error {
	pb, err := reload(snapshotId, playbookId)
	if err != nil {
		return err
	}
	// loads from snapshot table the  SnapshotId is 0, so we need fill it with
	// the updated snapshotId.
	pb.SnapshotId = snapshotId
	gCache.SetPlaybook(playbookId, pb)
	logx.L().Infof("update pb:%+v with snapshot id: %+v", playbookId, snapshotId)
	return nil
}

// GetWorkerTopic
func GetWorkerTopic(namespace string) string {
	namespace = fmt.Sprintf("Flow_%s", namespace)
	return namespace
}

// GetWorkerNamespace
func GetWorkerNamespace(topic string) string {
	x := strings.Split(topic, "Flow_")
	return x[1]
}
