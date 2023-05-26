/**
 * @Author: huaxr
 * @Description:
 * @File: proto
 * @Version: 1.0.0
 * @Date: 2021/9/7 下午3:35
 */

package appservice

type CreateAppReq struct {
	AppName     string `json:"app_name"`
	Description string `json:"description"`
	// 预估eps
	Eps int `json:"eps"`
}

type CreateAppResp struct {
	Id int `json:"id"`
}

type HandlerStatus string

const (
	Accept HandlerStatus = "accept"
	Reject HandlerStatus = "reject"
)

type UpdateAppInternalReq struct {
	Status  HandlerStatus `json:"status"`
	AppId   int           `json:"app_id"`
	Brokers string        `json:"brokers"`
	Eps     int           `json:"eps"`
}

type UpdateAppInternalResp struct {
	Id int `json:"id"`
}
