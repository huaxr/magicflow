// Author: huaxr
// Time:   2021/6/24 下午1:44
// Git:    huaxr

package request

import "github.com/huaxr/magicflow/core"

type HasAuthRes struct {
	Data map[string]interface{} `json:"data"`
}

type WorkerResponseReq struct {
	ServiceAddr string          `json:"service_addr"`
	Key         core.AcKey      `json:"key" form:"key" binding:"required"`
	Signature   *core.Signature `json:"signature"`
	Env         core.Env        `json:"env"`
	Output      interface{}     `json:"output" form:"output" binding:"required"`
	Token       string          `json:"token" form:"token" binding:"required"`
	HeartBeat   *core.HeartBeat `json:"heart_beat" form:"heart_beat"`
}

type WorkerExceptionReq struct {
	ServiceAddr string          `json:"service_addr"`
	Key         core.AcKey      `json:"key" form:"key" binding:"required"`
	Signature   *core.Signature `json:"signature"`
	Exception   interface{}     `json:"exception" form:"exception" binding:"required"`
	Token       string          `json:"token" form:"token" binding:"required"`
}

type WorkerAuth struct {
	Namespace string `json:"namespace" binding:"required"`
	Token     string `json:"token" binding:"required"`
}
