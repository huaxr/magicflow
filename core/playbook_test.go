// Author: huaxr
// Time:   2021/6/6 下午2:11
// Git:    huaxr

package core

import (
	"encoding/json"
	"testing"
)

var raw = `
{"id":1,"states":
{"a":{"branches":{"choice":[{"express":"1==1","next_node":"f","priority":1}]},"hook":false,"node_code":"a","node_type":"node","task":{"configuration":{},"input_attribute":"$$.","node_code":"a","playbook_id":138,"snapshot_id":1201,"status":"TaskReady","xrn":{"namespace":"","region":"rn","service":"proxy","task_info":"F","task_type":"normal"}}},
"b":{"branches":{"choice":[{"express":"$$$. == 10000","next_node":"g","priority":1},{"express":"1==1","next_node":"a","priority":1}],"parallel":[{"express":"","next_node":"d","priority":0},{"express":"","next_node":"c","priority":0}]},"hook":false,"node_code":"b","node_type":"node","task":{"configuration":{},"input_attribute":"$$.","node_code":"b","playbook_id":138,"snapshot_id":1201,"status":"TaskReady","xrn":{"namespace":"","region":"rn","service":"proxy","task_info":"B","task_type":"normal"}}},
"c":{"branches":{"parallel":[{"express":"","next_node":"f","priority":0}]},"hook":false,"node_code":"c","node_type":"node","task":{"configuration":{},"input_attribute":["$.e.nums[0:2]",88,99],"node_code":"c","playbook_id":138,"snapshot_id":1201,"status":"TaskReady","xrn":{"namespace":"","region":"rn","service":"proxy","task_info":"D","task_type":"normal"}}},
"d":{"branches":{"parallel":[{"express":"","next_node":"f","priority":0}]},"hook":false,"node_code":"d","node_type":"node","task":{"configuration":{},"input_attribute":{"A":"$.e.nums[0]","B":"$.b","X":[1,2,3],"Y":"$$$."},"node_code":"d","playbook_id":138,"snapshot_id":1201,"status":"TaskReady","xrn":{"namespace":"","region":"rn","service":"proxy","task_info":"C","task_type":"normal"}}},
"e":{"branches":{"parallel":[{"express":"","next_node":"b","priority":0}]},"hook":false,"node_code":"e","node_type":"node","task":{"configuration":{},"input_attribute":"$$.","node_code":"e","playbook_id":138,"snapshot_id":1201,"status":"TaskReady","xrn":{"namespace":"","region":"rn","service":"proxy","task_info":"A","task_type":"normal"}}},
"f":{"branches":{},"hook":false,"node_code":"f","node_type":"node","task":{"configuration":{},"input_attribute":"$$.","node_code":"f","playbook_id":138,"snapshot_id":1201,"status":"TaskReady","xrn":{"namespace":"","region":"rn","service":"proxy","task_info":"F","task_type":"normal"}}},
"g":{"branches":{"parallel":[{"express":"","next_node":"f","priority":0}]},"hook":false,"node_code":"g","node_type":"node","task":{"configuration":{},"input_attribute":"$$.","node_code":"g","playbook_id":138,"snapshot_id":1201,"status":"TaskReady","xrn":{"namespace":"","region":"rn","service":"proxy","task_info":"E","task_type":"normal"}}}},
"snapshot_id":1}`

func newTestPlaybook(t *testing.T) *playbook {
	newTestGlobal(t)
	var pb playbook
	err := json.Unmarshal([]byte(raw), &pb)
	if err != nil {
		t.Log(err)
		return nil
	}
	gCache.SetPlaybook(pb.Id, &pb)
	return &pb
}
