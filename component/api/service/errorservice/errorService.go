/**
 * @Author: huaxr
 * @Description:
 * @File: errorService
 * @Version: 1.0.0
 * @Date: 2021/9/10 上午11:04
 */

package errorservice

import (
	"context"
	"github.com/huaxr/magicflow/component/logx"
)

func NewError(ctx context.Context, req *NewErrorMsgReq) (NewErrorMsgResp, error) {
	logx.L().Errorf("NewError.Test", "text:%s", "这是一条测试的error提示")
	return NewErrorMsgResp{}, nil
}
