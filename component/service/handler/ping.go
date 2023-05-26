// Author: XinRui Hua
// Time:   2022/4/8 下午5:44
// Git:    huaxr

package handler

import (
	"context"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/service/proto"
)

type PingService struct{}

func (t *PingService) PingTest(ctx context.Context, req *proto.PingReq) (*proto.PongResponse, error) {
	logx.L().Infof("ping call %v", req.Message)
	return &proto.PongResponse{
		Status: 1,
	}, nil
}
