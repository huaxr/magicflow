package configController

import (
	"encoding/json"
	"github.com/huaxr/magicflow/component/api/auth"
	"github.com/huaxr/magicflow/component/api/controller"
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/dao/orm/models"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/service/client"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/huaxr/magicflow/pkg/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ConfigController struct {
	controller.BaseApiController
}

func (ctl *ConfigController) Router(router *gin.Engine) {
	// those two only for nsq cluster api
	// nsq 集群鉴权
	router.GET("/auth", ctl.ClientBrokerAuth)
	// nsq 集群回调
	router.POST("/notifycation", ctl.NsqAdminNotifycation)
	router.POST("/sql_shell", ctl.Sqlshell)

	// 返回rpc hosts 列表，用于prom监控的 自动 ip 变化监控
	router.GET("/hosts", ctl.Hosts)

	entry := router.Group("/config")
	// nsq lookup 列表（用于worker接入）
	entry.GET("/lookups", ctl.GetLookUps)
	// 用于worker接入时鉴权
	entry.POST("/auth", ctl.CheckAuth)
	// brokers 列表（用于管理员）
	entry.GET("/admin_brokers", ctl.GetBrokers)

}

type Object struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

func (ctl *ConfigController) Hosts(c *gin.Context) {
	var objs = make([]*Object, 0)

	var obj = new(Object)
	obj.Targets = client.GetHosts()
	obj.Labels = map[string]string{"label": "rpc"}

	objs = append(objs, obj)
	c.JSON(200, objs)
}

// only admin options will trigger notifycation
// when consumer create the topics, this notify will
// not be catched.
// with the advent of monitor process we should come up
// one way to this steps.
func (ctl *ConfigController) NsqAdminNotifycation(c *gin.Context) {
	jsons := make(map[string]interface{}) //注意该结构接受的内容
	if err := c.ShouldBind(&jsons); err != nil {
		logx.L().Errorf("NsqAdminNotifycation", "%v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	b, _ := json.Marshal(jsons)
	logx.L().Warnf("nsqadmin", string(b))
	c.JSON(http.StatusOK, "")
	return
}

func (ctl *ConfigController) Sqlshell(c *gin.Context) {
	jsons := make(map[string]interface{})
	if err := c.ShouldBind(&jsons); err != nil {
		logx.L().Errorf("NsqAdminNotifycation", "%v", err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}

	t, ok := jsons["token"]
	if !ok || t != "fakeadmin" {
		c.JSON(http.StatusOK, "token err")
		return
	}

	dao := orm.GetEngine()
	_, err := dao.Master().Exec(jsons["sql"])
	c.JSON(http.StatusOK, err)
	return

}

// get query by nsq broker
// "/auth?common_name=&remote_ip=10.73.35.27&secret=4BFE467B-FCBA-4519-BAC8-E9A3C57EDEB6&tls=false"
func (ctl *ConfigController) ClientBrokerAuth(c *gin.Context) {
	secretNamespace, _ := c.GetQuery("secret")
	ip, _ := c.GetQuery("remote_ip")

	logx.L().Infof("auth ip: %v, secret :%v", ip, secretNamespace)
	// is server
	// ip should allowed? no, local debug need. maybe online need this.
	if secretNamespace == confutil.GetConf().Queue.Nsq.Secret {
		var acc = make([]*request.AuthAccount, 0)

		acc = append(acc, &request.AuthAccount{
			Channels:    []string{".*"},
			Topic:       ".*",
			Permissions: []string{"subscribe", "publish"},
		})
		a := request.NsqAuth{
			Identity:       "hello xeflow server.",
			TTL:            3600 * 24 * 365,
			Authorizations: acc,
		}
		c.JSON(http.StatusOK, a)
		return
	}

	// is worker
	pair := strings.Split(secretNamespace, "?")
	if len(pair) != 2 {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "NOT_AUTHORIZED",
		})
		return
	}
	secret, domain := pair[0], pair[1]
	dao := orm.GetEngine()

	var ss = make([]models.App, 0)
	err := dao.Slave().Where("token = ? and app_name = ?", secret, domain).Find(&ss)
	if err != nil {
		logx.L().Errorf("BrokerAuth", "%v", err)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "NOT_AUTHORIZED",
		})
		return
	}

	if len(ss) != 1 {
		logx.L().Warnf("BrokerAuth", "NOT EXIST: %v, %v ip:%v", secret, domain, ip)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "NOT_AUTHORIZED",
		})
		return
	}

	var acc = make([]*request.AuthAccount, 0)

	acc = append(acc, &request.AuthAccount{
		Channels:    []string{".*"},
		Topic:       core.GetWorkerTopic(ss[0].AppName),
		Permissions: []string{"subscribe"},
	})
	a := request.NsqAuth{
		Identity:       "hello xeflow worker.",
		TTL:            3600 * 24 * 365,
		Authorizations: acc,
	}
	c.JSON(http.StatusOK, a)
	return
}

func (ctl *ConfigController) GetLookUps(c *gin.Context) {
	ctl.Success(c, confutil.GetConf().Queue.Nsq.Lookups)
}

func (ctl *ConfigController) GetBrokers(c *gin.Context) {
	uif, _ := c.Get("auth")
	if uif.(*auth.UserInfo).IsAdmin != 1 {
		ctl.Error(c, 200, "you are not super admin", nil)
		return
	}

	ctl.Success(c, strings.Split(confutil.GetConf().Queue.Nsq.Brokers, ","))
}

func (ctl *ConfigController) CheckAuth(c *gin.Context) {
	tag := "ConfigController.CheckAuth"
	req := &request.WorkerAuth{}
	if err := c.ShouldBind(req); err != nil {
		logx.L().Errorf(tag, "WorkerAuth bind error with param:[%v] err:[%v]", *req, err)
		ctl.Error(c, 200, err.Error(), nil)
		return
	}
	dao := orm.GetEngine()
	var apps = make([]models.App, 0)
	err := dao.Slave().Where("app_name = ? and token = ?", req.Namespace, req.Token).Find(&apps)
	if err != nil || len(apps) != 1 {
		ctl.Error(c, 200, "err", map[string]interface{}{"status": "no", "secret": "", "broker_size": 0})
		return
	}
	ctl.Success(c, map[string]interface{}{"status": "ok", "secret": apps[0].Token, "broker_size": len(strings.Split(apps[0].Brokers, ","))})
}
