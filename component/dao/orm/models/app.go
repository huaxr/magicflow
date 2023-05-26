package models

import (
	"time"
)

type App struct {
	Id            int       `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	AppName       string    `xorm:"not null default '' comment('app名称') unique VARCHAR(255)" json:"app_name"`
	User          string    `xorm:"not null default '0' comment('用户') VARCHAR(255)" json:"user"`
	Token         string    `xorm:"not null comment('是否有权限消费nsq') unique VARCHAR(255)" json:"token"`
	Brokers       string    `xorm:"not null comment('brokers队列') VARCHAR(255)" json:"brokers"`
	BrokerType    string    `xorm:"not null comment('broker类型') VARCHAR(255)" json:"broker_type"`
	Eps           int       `xorm:"not null comment('限流') INT(255)" json:"eps"`
	GroupId       string    `xorm:"not null comment('用户字段组') VARCHAR(11)" json:"group_id"`
	UpdateTime    time.Time `xorm:"not null comment('更新时间') DATETIME" json:"update_time"`
	Checked       int       `xorm:"not null default 0 comment('是否通过核审') TINYINT(1)" json:"checked"`
	Share         int       `xorm:"not null default 0 comment('是否共享') TINYINT(1)" json:"share"`
	Description   string    `xorm:"not null comment('描述') VARCHAR(255)" json:"description"`
	LastAliveTime time.Time `xorm:"not null comment('上次心跳时间') DATETIME" json:"last_alive_time"`
}
