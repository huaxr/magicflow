// Author: XinRui Hua
// Time:   2022/3/18 下午4:43
// Git:    huaxr

package service

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/huaxr/magicflow/component/plugin/limiter"

	"github.com/huaxr/magicflow/component/consensus"
	"github.com/huaxr/magicflow/component/helper/middleware"
	"github.com/huaxr/magicflow/component/logx"
	"github.com/huaxr/magicflow/component/service/handler"
	"github.com/huaxr/magicflow/component/service/proto"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/huaxr/magicflow/pkg/toolutil"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/wait"
)

func LaunchRpcServer(ctx context.Context, port string) (err error) {
	addr := toolutil.GetIp()
	key := fmt.Sprintf("%s/%s:%s", confutil.GetConf().Configuration.ServicesPrefix, addr, port)
	alive := consensus.KeepAlive{
		Key: key,
		Ttl: 8,
	}
	go wait.UntilWithContext(ctx, alive.Alive, time.Second*5)

	server := grpc.NewServer(
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcrecovery.StreamServerInterceptor(middleware.RecoveryInterceptor()),
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcrecovery.UnaryServerInterceptor(middleware.RecoveryInterceptor()),
		)),
	)

	proto.RegisterTriggerRpcServer(server, &handler.TriggerService{})
	proto.RegisterPlaybookRpcServer(server, &handler.PlaybookService{
		Name:  "playbook_rate",
		Limit: limiter.NewRateLimiter("", 5),
	})
	proto.RegisterAppRpcServer(server, &handler.AppService{})
	proto.RegisterConsoleRpcServer(server, &handler.ConsoleService{})
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logx.L().Panicf("LaunchRpcServer err %v", err)
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				logx.L().Warnf("rpc server shut down & revoke keepalive")
				server.GracefulStop()
				return
			}
		}
	}()

	return server.Serve(lis)
}
