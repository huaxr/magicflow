package models

import (
	"time"
)

type Tasks struct {
	Id            int       `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Name          string    `xorm:"not null default '' comment('函数名') VARCHAR(255)" json:"name"`
	Configuration string    `xorm:"not null default '' comment('函数的配置') VARCHAR(255)" json:"configuration"`
	Xrn           string    `xorm:"not null default '' comment('对应的xrn') VARCHAR(255)" json:"xrn"`
	Description   string    `xorm:"not null default '' comment('描述') VARCHAR(255)" json:"description"`
	AppId         int       `xorm:"not null comment('对应的appid') index INT(11)" json:"app_id"`
	Type          string    `xorm:"not null comment('任务类型') VARCHAR(255)" json:"type"`
	UpdateTime    time.Time `xorm:"not null comment('更新时间') DATETIME" json:"update_time"`
	InputExample  string    `xorm:"not null default '' comment('输入样例') VARCHAR(255)" json:"input_example"`
	OutputExample string    `xorm:"not null default '' comment('输出样例') VARCHAR(255)" json:"output_example"`
	User          string    `xorm:"not null comment('创建者') VARCHAR(255)" json:"user"`
}
