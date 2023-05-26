// Author: huaxr
// Time: 2022/7/1 6:19 下午
// Git: huaxr

package flags

import (
	"context"
	"fmt"
	"os"

	"github.com/huaxr/magicflow/component/service/proto"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/huaxr/magicflow/pkg/toolutil"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func Help(ctx *cli.Context) error {
	fmt.Println("help...")
	return nil
}

func GetCache(ctx *cli.Context) error {
	con, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%v", confutil.GetConf().Port.Service), grpc.WithInsecure())
	if err != nil {
		fmt.Printf("err when dial local:%v", err)
		return err
	}
	client := proto.NewConsoleRpcClient(con)
	resp, err := client.GetCache(context.Background(), &proto.GetCacheReq{})
	if err != nil {
		fmt.Printf("err when GetCache local:%v", err)
		return err
	}
	fmt.Printf("内存信息:\n%v", toolutil.Bytes2string(resp.Data))
	return nil
}

func GetAck(ctx *cli.Context) error {
	con, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%v", confutil.GetConf().Port.Service), grpc.WithInsecure())
	if err != nil {
		fmt.Printf("err when dial local:%v", err)
		return err
	}
	client := proto.NewConsoleRpcClient(con)
	resp, err := client.GetAck(context.Background(), &proto.GetAckReq{})
	if err != nil {
		fmt.Printf("err when GetCache local:%v", err)
		return err
	}
	fmt.Printf("ACK信息:\n%v", toolutil.Bytes2string(resp.Data))
	return nil
}

func GetTicker(ctx *cli.Context) error {
	con, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%v", confutil.GetConf().Port.Service), grpc.WithInsecure())
	if err != nil {
		fmt.Printf("err when dial local:%v", err)
		return err
	}
	client := proto.NewConsoleRpcClient(con)
	resp, err := client.GetTicker(context.Background(), &proto.GetTickerReq{})
	if err != nil {
		fmt.Printf("err when GetTickerReq local:%v", err)
		return err
	}
	fmt.Printf("定时任务信息:\n%v", toolutil.Bytes2string(resp.Data))
	return nil
}

func CleanTable(ctx *cli.Context) error {
	fmt.Printf("%v", os.Args)
	return nil
}
