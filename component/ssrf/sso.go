// Author: huaxr
// Time:   2022/1/14 上午10:22
// Git:    huaxr

package ssrf

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/huaxr/magicflow/component/logx"

	"github.com/huaxr/magicflow/pkg/request"
	"github.com/spf13/cast"
)

func (cli *HttpClient) GetTicket(appid, appkey string) (string, error) {
	var result request.TicketResult
	resp, err := cli.GetRestClient().R().
		Get(fmt.Sprintf("%s%s", cli.basicPath, fmt.Sprintf("?appid=%s&appkey=%s", appid, appkey)))

	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	return result.Ticket, nil
}

func (cli *HttpClient) GetUserInfo(token, ticket string) (request.UserInfoResult, error) {
	if len(token) == 0 || len(ticket) == 0 {
		return request.UserInfoResult{}, errors.New("token and ticket should not empty")
	}
	resp, err := cli.GetRestClient().R().
		Get(fmt.Sprintf("%s%s", cli.basicPath, fmt.Sprintf("?token=%s&ticket=%s", token, ticket)))
	if err != nil {
		logx.L().Errorf("GetUserInfo resp err %v", err)
		return request.UserInfoResult{}, err
	}

	var response map[string]interface{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		logx.L().Errorf("GetUserInfo.marshal err %v", err)
		return request.UserInfoResult{}, err
	}

	if r, ok := response["errcode"]; ok {
		if cast.ToInt(r) != 0 {
			return request.UserInfoResult{}, errors.New(cast.ToString(response["errmsg"]))
		}
	} else {
		return request.UserInfoResult{}, errors.New(fmt.Sprintf("format err:%+v", response))
	}

	var result request.UserInfoResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logx.L().Errorf("GetUserInfo.marshal err %v", err)
		return request.UserInfoResult{}, err
	}
	return result, nil
}

// https://doc-openapi.zhiyinlou.com/doc#/account/account
func (cli *HttpClient) GetMoreUserInfo(ticket, workcode string) (request.UserMoreInfoResult, error) {
	if len(ticket) == 0 || len(workcode) == 0 {
		return request.UserMoreInfoResult{}, errors.New("workcode and ticket should not empty")
	}
	resp, err := cli.GetRestClient().R().
		Get(fmt.Sprintf("%s%s", cli.basicPath, fmt.Sprintf("?ticket=%s&user_type=workcode&user_id=%s", ticket, workcode)))

	if err != nil {
		logx.L().Errorf("GetMoreUserInfo err %v", err)
		return request.UserMoreInfoResult{}, err
	}

	var response map[string]interface{}
	json.Unmarshal(resp.Body(), &response)
	if r, ok := response["errcode"]; ok {
		if cast.ToInt(r) != 0 {
			return request.UserMoreInfoResult{}, errors.New(cast.ToString(response["errmsg"]))
		}
	} else {
		return request.UserMoreInfoResult{}, errors.New(fmt.Sprintf("format err:%+v", response))
	}

	logx.L().Debugf("moreuserinfo %v", string(resp.Body()))
	var result request.UserMoreInfoResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		logx.L().Errorf("GetMoreUserInfo err %v", err)
		return request.UserMoreInfoResult{}, err
	}
	return result, nil
}
