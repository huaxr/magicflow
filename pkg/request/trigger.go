// Author: huaxr
// Time:   2022/1/14 上午10:34
// Git:    huaxr

package request

type TriggerPlaybook struct {
	// entey field changed
	AppName    string      `json:"app_name" form:"app_name" binding:"required"`
	PlaybookId int         `json:"playbook_id" form:"playbook_id" binding:"required"`
	AppToken   string      `json:"app_token" form:"app_token" binding:"required"`
	Sync       bool        `json:"sync"  form:"sync"`
	Data       interface{} `json:"data" form:"data" binding:"required"`
}

type HookStatePlaybook struct {
	// token?
	// hook specific node
	TraceId    uint   `json:"trace_id" form:"trace_id" binding:"required"`
	NodeCode   string `json:"node_code" form:"node_code" binding:"required"`
	SnapshotId int    `json:"snapshot_id" form:"snapshot_id" binding:"required"`
}
