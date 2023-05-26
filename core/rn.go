// Author: huaxr
// Time:   2021/6/15 下午4:08
// Git:    huaxr

package core

import (
	"fmt"

	"github.com/huaxr/magicflow/component/logx"

	"strings"

	"github.com/spf13/cast"
)

type TaskType string

const (
	// normal task points to node
	// local & remote task points to playbook
	NormalTask TaskType = "normal"
	NopTask    TaskType = "nop"

	LocalTaskSet  TaskType = "local"
	RemoteTaskSet TaskType = "remote"
)

// Xes Resource TaskInfo, for the cloud support functionality definition.
type xrn struct {
	// The service namespace that identifies the product (for example, faas, cronjob, proxy)
	Service string `json:"service"`
	// The region the resource resides in. Note that the xrn for some resources do not require a region, so this
	// component might be omitted.
	Region   string   `json:"region"`
	TaskType TaskType `json:"task_type"`
	// keep blank when use in current playbook, on the contrary, set the calling app's topic.
	Namespace string `json:"namespace"`
	// the resource name, like task_name or the calling playbook id.
	TaskInfo string `json:"task_info"`
}

func (xn *xrn) String() string {
	// e.g. proxy:cn:normal::ADD
	//      proxy:cn:nop::
	return fmt.Sprintf("magic:%s:%s:%s:%s:%s", xn.Service, xn.Region, xn.TaskType, xn.Namespace, xn.TaskInfo)
}

func NewNormalTaskRN(taskName string) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", "proxy", "rn", NormalTask, "", taskName)
}

func NewNopTaskRN(taskName string) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", "proxy", "rn", NopTask, "", taskName)
}

func NewLocalTaskRN(pbId int) string {
	return fmt.Sprintf("%s:%s:%s:%s:%d", "proxy", "cn", LocalTaskSet, "", pbId)
}

func NewRemoteTaskRN(remoteApp string, pbid int) string {
	return fmt.Sprintf("%s:%s:%s:%s:%d", "proxy", "cn", RemoteTaskSet, remoteApp, pbid)
}

func newMagicRn(rn string) *xrn {
	r := strings.Split(rn, ":")
	if len(r) != 5 {
		logx.L().Errorf("NewXesRn bad rn format.")
		return nil
	}

	return &xrn{
		Service:   r[0],
		Region:    r[1],
		TaskType:  TaskType(r[2]),
		Namespace: r[3],
		TaskInfo:  r[4],
	}
}

func (xn *xrn) getTaskName() string {
	if xn.TaskType != NormalTask {
		logx.L().Warnf("not node type has no task name")
		return ""
	}
	return xn.TaskInfo
}

func (xn *xrn) getCallId() int {
	if xn.TaskType != LocalTaskSet && xn.TaskType != RemoteTaskSet {
		logx.L().Warnf("not task set type has no call id")
		return -1
	}
	return cast.ToInt(xn.TaskInfo)
}

func (xn *xrn) getTaskType() TaskType {
	return xn.TaskType
}
