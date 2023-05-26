package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huaxr/magicflow/component/api"
	"github.com/huaxr/magicflow/component/dao/orm"
	"github.com/huaxr/magicflow/component/service/client"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
		time.Sleep(2 * time.Second)
		os.Exit(-1)
	}()
	orm.LaunchDbEngine()
	client.ServerFound()
	api.LaunchServer(ctx)
}
