package appController

import (
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/component/api/controller"
	"github.com/huaxr/magicflow/component/api/service/appservice"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"time"
)

type AppController struct {
	controller.BaseApiController
}

func (ctl *AppController) Router(router *gin.Engine) {
	entry := router.Group("/app")
	entry.GET("/get_info", ctl.GetAppInfo)

	// 获取用户有权限的app列表
	entry.GET("/list", ctl.GetAppList)

	// 获取app下所有剧本
	entry.GET("/get_playbook_list", ctl.GetPlaybookList)

	// 其它申请加入此app，用于租户场景
	entry.POST("/apply_use", ctl.ApplyUse)
	// 申请列表
	entry.GET("/get_apply_list", ctl.GetApplyList)
	// 通过申请
	entry.GET("/approve", ctl.Approve)

	// 数据库写入待核审的app，不走dcc集群化
	entry.POST("/create", ctl.CreateApp)
	//// 管理员审核通过app（触发dcc）
	//entry.POST("/admin_audit", ctl.UpdateAppInternal)

	// 管理员获取核审app列表
	entry.GET("/admin_get_audit_list", ctl.GetAuditList)
}

func (ctl *AppController) GetAppInfo(c *gin.Context) {
	appId, exist := c.GetQuery("id")
	if !exist || cast.ToInt(appId) <= 0 {
		logx.L().Errorf("id: %v not right!", appId)
		ctl.Error(c, 200, "err id", nil)
		return
	}

	dao := orm.GetEngine()

	var res = make([]models.App, 0)
	err := dao.Slave().Where("id = ?", appId).Cols("app_name", "user", "token", "eps").Find(&res)
	if err != nil {
		logx.L().Errorf("query app by id err: %+v", err)
		ctl.Error(c, 200, "query err", nil)
		return
	}

	ctl.Success(c, res)
}

func (ctl *AppController) GetAppList(c *gin.Context) {
	dao := orm.GetEngine()

	var res = make([]models.App, 0)
	var err error

	uif, _ := c.Get("auth")
	name := uif.(*auth.UserInfo).Account

	// 获取属于个人 或者 共享的app
	err = dao.Slave().Where("user = ? or share = 1", name).Find(&res)
	if err != nil {
		ctl.Error(c, 200, "query err", err)
		return
	}

	ctl.Success(c, res)
}

func (ctl *AppController) GetApplyList(c *gin.Context) {
	appId, exist := c.GetQuery("app_id")
	if !exist || cast.ToInt(appId) <= 0 {
		ctl.Error(c, 200, "err appId", nil)
		return
	}

	uif, _ := c.Get("auth")
	name := uif.(*auth.UserInfo).Account

	dao := orm.GetEngine()

	var apps = make([]models.App, 0)
	err := dao.Slave().Where("id = ? and user = ?", appId, name).Find(&apps)
	if err != nil {
		ctl.Error(c, 200, "query err", nil)
		return
	}

	if len(apps) == 0 {
		ctl.Error(c, 200, "you are not allowed to access this app", nil)
		return
	}

	var appusers = make([]models.AppUser, 0)
	dao.Master().Where("app_id = ?", appId).OrderBy("id desc").Find(&appusers)
	ctl.Success(c, appusers)
}

func (ctl *AppController) Approve(c *gin.Context) {
	appUserId, exist := c.GetQuery("appuser_id")
	if !exist || cast.ToInt(appUserId) <= 0 {
		ctl.Error(c, 200, "err appUserId", nil)
		return
	}
	dao := orm.GetEngine()

	// 查找对应的项
	var appusers = make([]models.AppUser, 0)
	dao.Slave().Where("id = ?", appUserId).Find(&appusers)
	if len(appusers) != 1 {
		ctl.Error(c, 200, "query not exist", nil)
		return
	}

	uif, _ := c.Get("auth")
	name := uif.(*auth.UserInfo).Account

	// 鉴权
	var apps = make([]models.App, 0)
	dao.Slave().Where("id = ? and user = ?", appusers[0].AppId, name).Find(&apps)
	if len(apps) == 0 {
		ctl.Error(c, 200, "you are not allowed to access this app", nil)
		return
	}

	// 鉴权成功，则更改标志位
	updatesappuser := models.AppUser{
		Checked: 1,
	}
	dao.Master().MustCols("checked").Where("id = ?", appUserId).Update(&updatesappuser)
	ctl.Success(c, "approve success")
}

func (ctl *AppController) GetPlaybookList(c *gin.Context) {
	appId, exist := c.GetQuery("app_id")
	if !exist || cast.ToInt(appId) <= 0 {
		logx.L().Errorf("appId: %v not right!", appId)
		ctl.Error(c, 200, "err appId", nil)
		return
	}
	uif, _ := c.Get("auth")
	user := uif.(*auth.UserInfo)
	err := user.HasAppPermit(cast.ToInt(appId), true)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	dao := orm.GetEngine()

	var res = make([]models.Playbook, 0)
	dao.Slave().Where("app_id = ?", appId).OrderBy("id desc").Find(&res)
	ctl.Success(c, res)
}

func (ctl *AppController) CreateApp(c *gin.Context) {
	var params appservice.CreateAppReq
	err := c.ShouldBind(&params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ret, err := appservice.CreateApp(c, &params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ctl.Success(c, ret)
	return
}

func (ctl *AppController) ApplyUse(c *gin.Context) {
	appId, exist := c.GetQuery("app_id")
	if !exist || cast.ToInt(appId) <= 0 {
		ctl.Error(c, 200, "appid err", nil)
		return
	}

	uif, _ := c.Get("auth")
	uid := uif.(*auth.UserInfo).Uid

	dao := orm.GetEngine()

	// Check the uniqueness of the app_name
	var res = make([]models.App, 0)
	err := dao.Slave().Where("id = ?", appId).Find(&res)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	// Create new app
	create := models.AppUser{
		UserId:     uid,
		AppId:      cast.ToInt(appId),
		Checked:    0,
		CreateTime: time.Now(),
	}
	_, err = dao.Master().Insert(&create)

	ctl.Success(c, "apply success, wait audit")
}

func (ctl *AppController) UpdateAppInternal(c *gin.Context) {
	var params appservice.UpdateAppInternalReq
	err := c.ShouldBind(&params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	uif, _ := c.Get("auth")
	if uif.(*auth.UserInfo).IsAdmin != 1 {
		ctl.Error(c, 200, "you are not super admin", nil)
		return
	}

	ret, err := appservice.UpdateAppInternal(c, &params)
	if err != nil {
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	ctl.Success(c, ret)
	return
}

func (ctl *AppController) GetAuditList(c *gin.Context) {
	uif, _ := c.Get("auth")
	if uif.(*auth.UserInfo).IsAdmin != 1 {
		ctl.Error(c, 200, "you are not super admin", nil)
		return
	}
	status := c.DefaultQuery("status", "")
	dao := orm.GetEngine()

	// 已审核 （和提交接口参数不一致）
	if status == "checked" {
		var res = make([]models.App, 0)
		err := dao.Slave().Where("checked = ? or checked = ?", 1, -1).Find(&res)
		if err != nil {
			ctl.Error(c, 200, err.Error(), nil)
			return
		}
		ctl.Success(c, res)
		return
	} else {
		// 未审核
		var res = make([]models.App, 0)
		err := dao.Slave().Where("checked = ?", 0).Find(&res)
		if err != nil {
			ctl.Error(c, 200, err.Error(), nil)
			return
		}
		ctl.Success(c, res)
		return
	}
}
