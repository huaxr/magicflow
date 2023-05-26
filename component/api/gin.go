// Author: huaxr
// Time:   2021/6/8 下午3:52
// Git:    huaxr

package api

import (
	"context"
	"fmt"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/pkg/confutil"
	"io/ioutil"

	"github.com/huaxr/magicflow/component/api/controller/appController"
	"github.com/huaxr/magicflow/component/api/controller/authController"
	"github.com/huaxr/magicflow/component/api/controller/configController"
	"github.com/huaxr/magicflow/component/api/controller/playbookcontroller"
	"github.com/huaxr/magicflow/component/api/controller/taskController"
	"github.com/huaxr/magicflow/component/api/controller/triggercontroller"
	"github.com/huaxr/magicflow/component/api/middleware/normal"
	"github.com/huaxr/magicflow/pkg/accutil"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
	"sync"
)

type Router struct {
	Router *gin.Engine
	once   *sync.Once
	mod    string
}

func init() {
	if runtime.NumCPU() >= 8 {
		runtime.GOMAXPROCS(8)
	}
}

func LaunchServer(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			logx.L().Errorf("LaunchServerv panic %v", err)
		}
	}()
	e := new(Router)
	e.Router = gin.Default()
	e.Router.LoadHTMLGlob("static/templates/*")
	e.Router.StaticFS("/static", http.Dir("static"))

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	e.once = new(sync.Once)
	e.registerMiddleware(normal.LoginRequired())
	e.registerRouter(e.Router)
	go accutil.Thread("WEB-SERVER", func() {
		_ = e.Router.Run(fmt.Sprintf(":%s", confutil.GetConf().Port.Api))
	})

	select {
	case <-ctx.Done():
		logx.L().Errorf("main thread exit.")
	}
}

func (r *Router) registerMiddleware(middleware ...gin.HandlerFunc) {
	// By default gin.DefaultWriter = os.Stdout
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.once.Do(func() {
		for _, m := range middleware {
			r.Router.Use(m)
		}
	})
}

func (r *Router) registerRouter(router *gin.Engine) {
	new(playbookcontroller.PlayBookController).Router(router)
	new(configController.ConfigController).Router(router)
	new(appController.AppController).Router(router)
	new(taskController.TaskController).Router(router)
	new(authController.AuthController).Router(router)

	// 以下未rpc重构的router
	// 所有涉及 dcc G 等操作变更的均需要重写
	new(triggercontroller.RpcController).Router(router)
	new(playbookcontroller.RpcController).Router(router)
	new(appController.RpcController).Router(router)

}
