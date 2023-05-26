package models

import (
	"time"
)

type Playbook struct {
	Id          int       `xorm:"not null pk autoincr comment('主键') INT(11)" json:"id"`
	AppId       int       `xorm:"not null comment('对应的app id') INT(11)" json:"app_id"`
	SnapshotId  int       `xorm:"not null comment('快照版本id，用于动态切换') INT(11)" json:"snapshot_id"`
	User        string    `xorm:"not null comment('创建人') VARCHAR(255)" json:"user"`
	Name        string    `xorm:"not null comment('剧本名称') VARCHAR(255)" json:"name"`
	Enable      int       `xorm:"not null default 0 comment('是否开启') TINYINT(255)" json:"enable"`
	Description string    `xorm:"not null comment('剧本描述') VARCHAR(255)" json:"description"`
	Token       string    `xorm:"not null comment('授权token，用于被远程调用时鉴权') VARCHAR(255)" json:"token"`
	UpdateTime  time.Time `xorm:"not null comment('更新时间') DATETIME" json:"update_time"`
}
