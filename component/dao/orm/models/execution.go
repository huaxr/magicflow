package models

type Execution struct {
	Id         int    `xorm:"not null pk autoincr comment('主键') INT(11)" json:"id"`
	TraceId    string `xorm:"comment('执行链路id') index VARCHAR(255)" json:"trace_id"`
	Sequence   int    `xorm:"comment('节点code') INT(11)" json:"sequence"`
	NodeCode   string `xorm:"comment('节点id') index VARCHAR(32)" json:"node_code"`
	Domain     string `xorm:"comment('所属域') VARCHAR(255)" json:"domain"`
	Status     string `xorm:"comment('状态') VARCHAR(30)" json:"status"`
	PlaybookId int    `xorm:"not null comment('剧本id') INT(11)" json:"playbook_id"`
	Extra      string `xorm:"comment('额外信息') TEXT" json:"extra"`
	Timestamp  string `xorm:"comment('时间戳') VARCHAR(55)" json:"timestamp"`
	SnapshotId int    `xorm:"comment('快照id') index INT(11)" json:"snapshot_id"`
	Chain      string `xorm:"comment('执行链路') index VARCHAR(255)" json:"chain"`
}
