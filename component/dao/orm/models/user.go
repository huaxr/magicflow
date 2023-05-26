package models

import (
	"time"
)

type User struct {
	Id           int       `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Account      string    `xorm:"not null comment('用户') VARCHAR(255)" json:"account"`
	Name         string    `xorm:"not null comment('邮箱') VARCHAR(255)" json:"name"`
	Workcode     string    `xorm:"not null comment('工号') VARCHAR(255)" json:"workcode"`
	CreateTime   time.Time `xorm:"not null comment('创建时间') DATETIME" json:"create_time"`
	DeptId       string    `xorm:"not null comment('部门') VARCHAR(255)" json:"dept_id"`
	DeptName     string    `xorm:"not null comment('部门') VARCHAR(255)" json:"dept_name"`
	DeptFullName string    `xorm:"not null comment('部门') VARCHAR(255)" json:"dept_full_name"`
	Email        string    `xorm:"not null comment('email') VARCHAR(255)" json:"email"`
	Avatar       string    `xorm:"not null comment('头像') VARCHAR(255)" json:"avatar"`
	SuperAdmin   int       `xorm:"not null default 0 comment('是否是超级管理员') TINYINT(1)" json:"super_admin"`
}
