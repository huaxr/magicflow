// Author: XinRui Hua
// Time:   2022/3/22 下午2:26
// Git:    huaxr

package test

import (
	"context"

	"github.com/huaxr/magicflow/component/logx"

	"net"
	"testing"

	"google.golang.org/grpc"
)

type HandlerService struct{}

func (t *HandlerService) Test(ctx context.Context, req *TestReq) (*TestResponse, error) {
	logx.L().Infof("call me success %v", req.Payload)
	return &TestResponse{
		Data: "hello！ " + req.Payload,
	}, nil
}

func TestRpc(t *testing.T) {
	conn, err := grpc.Dial("10.74.152.189:8888", grpc.WithInsecure())
	if err != nil {
		t.Log(err)
		return
	}

	defer conn.Close()

	client := NewTestRpcClient(conn)
	resp, err := client.Test(context.Background(), &TestReq{
		Payload: "tmp",
	})
	if err != nil {
		t.Log(err)
		return
	}

	t.Log(resp.Data)

}

func TestServer(t *testing.T) {
	server := grpc.NewServer()
	RegisterTestRpcServer(server, &HandlerService{})
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		t.Log(err)
	}
	server.Serve(lis)
}
