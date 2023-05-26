// Author: huaxr
// Time:   2021/12/16 下午4:21
// Git:    huaxr

package taskController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/component/api/controller"
	"github.com/huaxr/magicflow/component/api/service/taskService"
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/spf13/cast"
)

type TaskController struct {
	controller.BaseApiController
}

func (ctl *TaskController) Router(router *gin.Engine) {
	entry := router.Group("/task")
	// 获取当前app下得任务列表
	entry.GET("/get_tasks", ctl.GetTasksList)
	entry.GET("/get_task_detail", ctl.GetTaskDetail)
	// 创建普通任务
	entry.POST("/create_normal", ctl.CreateNormalTask)
	entry.POST("/create_local", ctl.CreateLocalTask)
	// 创建远程任务
	entry.POST("/create_remote", ctl.CreateRemoteTask)
	// task uodate should be condi
}

func (ctl *TaskController) GetTasksList(c *gin.Context) {
	appId, exist := c.GetQuery("app_id")
	if !exist || cast.ToInt(appId) <= 0 {
		ctl.Error(c, 200, "err appId", nil)
		return
	}
	dao := orm.GetEngine()

	// appid 要存在
	var apps = make([]models.App, 0)
	dao.Slave().Where("id = ?", appId).Find(&apps)
	if len(apps) == 0 {
		ctl.Error(c, 200, fmt.Sprintf("appid %s not exist", appId), nil)
		return
	}

	// app必须经过核审
	if apps[0].Checked != 1 {
		ctl.Error(c, 200, fmt.Sprintf("app %s not checked yet", appId), nil)
		return
	}

	// 鉴权
	uif, _ := c.Get("auth")
	user := uif.(*auth.UserInfo)
	err := user.HasAppPermit(cast.ToInt(appId), true)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
	}

	var res = make([]models.Tasks, 0)
	err = dao.Slave().Where("app_id = ?", appId).OrderBy("id desc").Find(&res)
	if err != nil {
		ctl.Error(c, 200, "query err", nil)
		return
	}

	ctl.Success(c, res)
}

func (ctl *TaskController) GetTaskDetail(c *gin.Context) {
	taskId, exist := c.GetQuery("id")
	if !exist || cast.ToInt(taskId) <= 0 {
		ctl.Error(c, 200, "err id", nil)
		return
	}

	dao := orm.GetEngine()
	var res = make([]models.Tasks, 0)
	err := dao.Slave().Where("id = ?", taskId).Find(&res)
	if err != nil {
		ctl.Error(c, 200, "query err", nil)
		return
	}
	if len(res) != 1 {
		ctl.Success(c, gin.H{})
		return
	}
	ctl.Success(c, res[0])
}

func (ctl *TaskController) CreateNormalTask(c *gin.Context) {
	var params taskService.CreateNormalTaskReq
	err := c.ShouldBind(&params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ret, err := taskService.CreateTask(c, &params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ctl.Success(c, ret)
	return
}

func (ctl *TaskController) CreateLocalTask(c *gin.Context) {
	var params taskService.CreateLocalTaskReq
	err := c.ShouldBind(&params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	// no end nodes could not used by localtask
	ret, err := taskService.CreateLocalTask(c, &params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ctl.Success(c, ret)
	return
}

func (ctl *TaskController) CreateRemoteTask(c *gin.Context) {
	var params taskService.CreateRemotePlaybookTaskReq
	err := c.ShouldBind(&params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ret, err := taskService.CreateRemoteTask(c, &params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ctl.Success(c, ret)
	return
}
