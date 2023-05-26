package core

import (
	"encoding/json"
	"strings"

	"github.com/huaxr/magicflow/pkg/toolutil"
)

type (
	TaskStatus string
)

const (
	// when config playbook, Ready will be set.
	Ready TaskStatus = "TaskReady"
	// when trigger called.
	Start TaskStatus = "TaskStart"
	// when ToServer flying in the local channel.
	Flying   TaskStatus = "TaskFlying"
	Executed TaskStatus = "TaskExecuted"
	Hooked   TaskStatus = "TaskHooked"
	Fail     TaskStatus = "TaskFail"
)

func newAttribute(attrStr string) (interface{}, error) {
	attrStr = strings.Trim(attrStr, " ")
	if strings.HasPrefix(attrStr, "{") || strings.HasPrefix(attrStr, "[") {
		var att interface{}
		err := json.Unmarshal(toolutil.String2Byte(attrStr), &att)
		if err != nil {
			return nil, err
		}
		return att, nil
	}
	return attrStr, nil
}

// task shared by local pb's node and transform by message
// there are redundant fields should be filled.
type task struct {
	// redundant fields keeps going in message body
	Status     TaskStatus `json:"status"`
	PlaybookId int        `json:"playbook_id"`
	SnapshotId int        `json:"snapshot_id"`
	NodeCode   string     `json:"node_code"`
	// input restriction relations for fill the input body.
	Input  interface{} `json:"input,omitempty"`
	Output interface{} `json:"output,omitempty"`
	// retrospect the service deploy requirements delays the
	// whole functionalities tmp, task should provide sandbox
	// environment to enable tmp regression integrity.
	Sandbox interface{} `json:"sandbox,omitempty"`

	// task sharing Meta info, which binding to a node
	// encapsulation of the input
	// attribute e.g.  {"sub_number": "$.ADD.add_number"} or {"sub_number": "$$.add_number"}
	// in order to support intrinsic values, like {"k": 123}, the Map should support interface{} parameters.
	// there are four occasions:
	//   1: Trigger $$$.
	//   2: Last    $$.
	//	 3: SpecificNode $.a
	//   4: raw plain data
	//   5. Ctx data. $._env
	// configuration should in obedience to those occasions.
	InputAttribute interface{} `json:"input_attribute,omitempty"`

	// OutPut should defined here
	//OutputAttribute interface{} `json:"output_attribute,omitempty"`

	Xrn *xrn `json:"xrn"`
	// some configuration binding to this node which give instructions toSlave.
	Configuration *configuration `json:"configuration,omitempty"`
}

func (t *task) getStatus() TaskStatus {
	return t.Status
}

func (t *task) GetConfiguration() *configuration {
	if t.Configuration == nil {
		return &configuration{
			Retry:         0,
			Timeout:       0,
			BeforeExecute: nil,
			AfterExecute:  nil,
			Token:         "",
		}
	}
	return t.Configuration
}

func (t *task) getPlaybookId() int { return t.PlaybookId }

func (t *task) getSnapshotId() int { return t.SnapshotId }

func (t *task) getNodeCode() string { return t.NodeCode }

func (t *task) getInput() interface{} { return t.Input }

func (t *task) getOutput() interface{} { return t.Output }

func (t *task) setStatus(state TaskStatus) { t.Status = state }

func (t *task) setOutput(out interface{}) { t.Output = out }

func (t *task) setInput(in interface{}) { t.Input = in }

func (t *task) getXrn() *xrn {
	if t.Xrn == nil {
		return &xrn{}
	}
	return t.Xrn
}

// task always binding to one node which tells this option is
// hook or not, if hook set, then record the message.
func (t *task) isHook() bool {
	p, ok := gCache.GetPlaybook(t.PlaybookId)
	if !ok {
		return false
	}
	// when we switch snapshot, remove some node or rollback to a
	// node limited playbook, which cause panic because n would be
	// not exists.
	n, ok := p.GetNode(t.NodeCode)
	if n == nil {
		return false
	}
	return n.isHook()
}
