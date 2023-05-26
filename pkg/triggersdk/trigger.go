// Author: huaxr
// Time:   2021/6/22 下午1:58
// Git:    huaxr

package triggersdk

import (
	"fmt"
	"github.com/huaxr/magicflow/pkg/request"
)

// TriggerStateMachine Trigger a stateMachine. You can start a execution of specified stateMachine by passing stateMachineID in data.
// Also, it's available to trigger an unregistered stateMachine by passing it's definition in data.
func (cli *HttpClient) ExecutePlaybook(req *request.TriggerPlaybook) (interface{}, error) {
	var result map[string]interface{}
	resp, err := cli.restClient.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", req.AppToken)).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&result).
		Post(fmt.Sprintf("%s%s", cli.basicPath, "/trigger/execute"))

	if err != nil {
		return nil, err
	}
	objResp := resp.Result()
	return objResp, nil
}

func (cli *HttpClient) HookPlaybook(req *request.HookStatePlaybook) (interface{}, error) {
	var result map[string]interface{}
	resp, err := cli.restClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&result).
		Post(fmt.Sprintf("%s%s", cli.basicPath, "/trigger/hook"))

	if err != nil {
		return nil, err
	}
	objResp := resp.Result()
	return objResp, nil
}
