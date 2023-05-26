package models

type Node struct {
	Id int `xorm:"not null pk autoincr comment('主键') INT(11)"`
	//NodeName       string `xorm:"comment('节点中文名称') VARCHAR(255)"`
	NodeCode       string `xorm:"not null comment('节点代号') VARCHAR(255)"`
	InputAttribute string `xorm:"comment('输入属性') VARCHAR(255)"`
	PlaybookId     int    `xorm:"not null comment('剧本id') INT(255)"`
	Dependency     string `xorm:"comment('多个入度的父亲节点，只有一个时忽略不计') VARCHAR(255)"`
	TaskId         int    `xorm:"not null comment('绑定的任务集') INT(11)"`
	Hook           int    `xorm:"default 0 comment('是否hook') TINYINT(1)"`
}
