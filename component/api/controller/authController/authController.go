// Author: huaxr
// Time:   2022/1/5 下午6:30
// Git:    huaxr

package authController

import (
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/gin-gonic/gin"
	"github.com/huaxr/magicflow/component/api/controller"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/ssrf"
	"github.com/huaxr/magicflow/pkg/accutil"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/huaxr/magicflow/pkg/jwtutil"
	"strings"
	"time"
)

type AuthController struct {
	controller.BaseApiController
}

var (
	getticket = ssrf.NewHttpClient("https://xxx/basic/get_ticket")
	// limited information
	verify     = ssrf.NewHttpClient("https://xxx/api/v1/sso/verify")
	verifymore = ssrf.NewHttpClient("https://xxx/cmpts/data/account/v2/user/get")
)

func getTicket() string {
	// https://sso.zhiyinlou.com/portal/login/1541377518
	appid, appkey := confutil.GetConf().Configuration.Appid, confutil.GetConf().Configuration.Appkey
	ticket, err := getticket.GetTicket(appid, appkey)
	if err != nil {
		logx.L().Errorf("getTicket err %v", err)
		return ""
	}
	return ticket
}

func (ctl *AuthController) Router(router *gin.Engine) {
	// sso callback
	router.GET("/callback", ctl.Callback)
}

func (ctl *AuthController) Callback(c *gin.Context) {
	token, exist := c.GetQuery("token")
	if !exist || len(token) == 0 {
		ctl.Error(c, 200, "err token", nil)
		return
	}
	redirect, exist := c.GetQuery("redirect")
	if !exist || len(redirect) == 0 {
		ctl.Error(c, 200, "err redirect", nil)
		return
	}
	res, err := verify.GetUserInfo(token, getTicket())
	if err != nil {
		logx.L().Errorf("Callback.verify err %v", err)
		ctl.Error(c, 200, "err", err)
		return
	}

	dao := orm.GetEngine()

	var users = make([]models.User, 0)
	var user models.User
	err = dao.Slave().Where("account = ?", res.Data.Account).Find(&users)
	if err != nil {
		logx.L().Errorf("Callback.user err %v", err)
		ctl.Error(c, 200, "err", err)
		return
	}

	if len(users) == 0 {
		more, err := verifymore.GetMoreUserInfo(getTicket(), res.Data.Workcode)
		if err != nil {
			logx.L().Errorf("Callback moreverify %v", err)
			ctl.Error(c, 200, err.Error(), nil)
			return
		}

		var deptid, deptname, deptfullname string
		if len(more.Data.DeptInfo) > 0 {
			deptid = more.Data.DeptInfo[0].EhrDeptId
			deptname = more.Data.DeptInfo[0].DeptName
			deptfullname = more.Data.DeptInfo[0].DeptFullName
		}

		var superuser int
		lds := strings.Split(confutil.GetConf().Configuration.Superadmin, ",")
		if accutil.ContainsStr(lds, res.Data.Account) {
			superuser = 1
		}
		user = models.User{
			Account:      res.Data.Account,
			Name:         res.Data.Name,
			Workcode:     res.Data.Workcode,
			DeptId:       deptid,
			DeptName:     deptname,
			DeptFullName: deptfullname,
			Email:        more.Data.Email,
			Avatar:       more.Data.Avatar,
			CreateTime:   time.Now(),
			SuperAdmin:   superuser,
		}
		dao.Master().Insert(&user)
	} else {
		user = users[0]
	}

	var userInfo = map[string]interface{}{
		"uid":           user.Id,
		"account":       user.Account,
		"name":          user.Name,
		"avatar":        user.Avatar,
		"deptname":      user.DeptName,
		"deptfulltname": user.DeptFullName,
		"is_admin":      user.SuperAdmin,
	}

	t, err := jwtutil.GenTokenString(userInfo)
	if err != nil {
		logx.L().Errorf("Callback jwt %v", err)
		ctl.Error(c, 200, "err", err)
		return
	}

	c.SetCookie("magic", t, 24*3600*365, "", c.Request.Header.Get("Referer"), false, false)
	c.Redirect(302, redirect)
}
