// Author: huaxr
// Time:   2022/1/6 下午5:04
// Git:    huaxr

package auth

import (
	"errors"
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/dao/orm/models"
)

type UserInfo struct {
	// db id
	Uid int
	// huaxinrui
	Account string
	// 华心瑞
	Name         string
	Avatar       string
	DeptName     string
	DeptFullName string
	// 是否是管理员
	IsAdmin int
}

func (user *UserInfo) HasAppPermit(appid int, judgeshare bool) error {
	dao := orm.GetEngine()
	// 鉴权是否有权限提交该剧本id
	// app必须经过核审
	var apps = make([]models.App, 0)
	err := dao.Slave().Where("id = ? and checked = 1", appid).Find(&apps)
	if len(apps) != 1 || err != nil {
		return errors.New("illegal app")
	}
	if judgeshare && apps[0].Share == 1 {
		// 共享app所有人都可以切换快照
		return nil
	} else {
		// 如果当前用户不是创建者
		if apps[0].User != user.Account {
			// 判断是否申请过该应用的权限
			var appusers = make([]models.AppUser, 0)
			err = dao.Slave().Where("app_id = ? and user_id = ? and checked = 1", appid, user.Uid).Find(&appusers)
			if len(appusers) == 0 || err != nil {
				return errors.New("you are not allowed to access this page")
			}
		}
	}
	return nil
}
