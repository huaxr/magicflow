// Author: huaxr
// Time: 2022/7/1 6:43 下午
// Git: huaxr

package handler

import (
	"context"

	"github.com/huaxr/magicflow/component/ticker"

	"github.com/huaxr/magicflow/component/service/proto"
	"github.com/huaxr/magicflow/core"
)

type ConsoleService struct{}

func (t *ConsoleService) GetCache(ctx context.Context, req *proto.GetCacheReq) (*proto.GetCacheResponse, error) {
	b := core.GetExporter().Export("cache")

	return &proto.GetCacheResponse{
		Data: b,
	}, nil
}

func (t *ConsoleService) GetAck(ctx context.Context, req *proto.GetAckReq) (*proto.GetAckResponse, error) {
	b := core.GetExporter().Export("ack")
	return &proto.GetAckResponse{
		Data: b,
	}, nil
}

func (t *ConsoleService) GetTicker(ctx context.Context, req *proto.GetTickerReq) (*proto.GetTickerResponse, error) {
	b := ticker.GetManager().Export("")
	return &proto.GetTickerResponse{
		Data: b,
	}, nil
}
