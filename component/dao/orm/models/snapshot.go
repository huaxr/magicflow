package models

import (
	"time"
)

type Snapshot struct {
	Id         int       `xorm:"not null pk autoincr comment('主键') INT(11)" json:"id"`
	PlaybookId int       `xorm:"not null comment('对应的剧本id') index INT(11)" json:"playbook_id"`
	Snapshot   string    `xorm:"not null comment('剧本快照body体') TEXT" json:"snapshot"`
	Rawbody    string    `xorm:"not null comment('前端传来的元数据') TEXT" json:"rawbody"`
	Checksum   string    `xorm:"not null comment('校验和') index VARCHAR(255)" json:"checksum"`
	UpdateTime time.Time `xorm:"not null comment('更新时间') DATETIME" json:"update_time"`
	AppId      int       `xorm:"not null comment('app') index INT(11)" json:"app_id"`
	Snapname   string    `xorm:"not null comment('快照名称') VARCHAR(255)" json:"snapname"`
	User       string    `xorm:"not null comment('创建人') VARCHAR(255)" json:"user"`
}
