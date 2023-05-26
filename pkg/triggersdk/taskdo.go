// Author: huaxr
// Time:   2021/6/24 下午2:58
// Git:    huaxr

package triggersdk

import (
	"fmt"
	"github.com/huaxr/magicflow/pkg/request"
)

func (cli *HttpClient) ReportTask(req *request.WorkerResponseReq) (interface{}, error) {
	var result map[string]interface{}
	resp, err := cli.restClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", req.Token)).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&result).
		Post(fmt.Sprintf("%s%s", cli.basicPath, "/trigger/worker_response"))

	if err != nil {
		return nil, err
	}
	objResp := resp.Result()
	return objResp, nil
}

func (cli *HttpClient) ReportException(req *request.WorkerExceptionReq) (interface{}, error) {
	var result map[string]interface{}
	resp, err := cli.restClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", req.Token)).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&result).
		Post(fmt.Sprintf("%s%s", cli.basicPath, "/trigger/worker_exception"))

	if err != nil {
		return nil, err
	}
	objResp := resp.Result()
	return objResp, nil
}
