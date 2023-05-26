package models

type Edge struct {
	Id         int    `xorm:"not null pk autoincr comment('主键') INT(11)"`
	PlaybookId int    `xorm:"not null comment('剧本底下的所有边') INT(11)"`
	FromNode   string `xorm:"not null comment('边的起点') VARCHAR(255)"`
	ToNode     string `xorm:"not null comment('边的终点') VARCHAR(255)"`
	Condition  string `xorm:"not null comment('bool表达式') VARCHAR(255)"`
	Priority   int    `xorm:"not null comment('优先级') INT(255)"`
}
