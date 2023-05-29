// Author:
// Time:   2021/6/7 下午8:45
// Git:    huaxr

package main

import (
	ctx2 "context"
	"encoding/json"
	"log"
	"time"

	"github.com/huaxr/magicflow/component/ticker"

	"github.com/huaxr/magicflow/pkg/request"
	"github.com/huaxr/magicflow/pkg/triggersdk"
)

var cli = triggersdk.NewHttpClient("http://127.0.0.1:8080")

//var cli = triggersdk.NewHttpClient("http://api.xueersi.com/orchestration/")

func hook() {
	res, err := cli.HookPlaybook(&request.HookStatePlaybook{
		TraceId:    1543857526886117385,
		NodeCode:   "62d0d7",
		SnapshotId: 1209,
	})
	log.Println(res, err)
}

func loop() {
	f := func() {
		res, _ := cli.ExecutePlaybook(&request.TriggerPlaybook{
			AppToken:   "dXpuYTVjNzZpYWNuaThx",
			AppName:    "first_example",
			PlaybookId: 138,
			Data:       1000,
			Sync:       true,
		})

		b, err := json.Marshal(res)
		if err != nil {
			log.Println("err", err)
			return
		}
		log.Println(string(b))
	}
	tick := time.NewTicker(1000 * time.Microsecond)
	job := ticker.NewJob(ctx2.Background(), "loop", tick, f)
	ticker.GetManager().Register(job)

}

// bug 应该是多分支的问题，存在产生两次的情况
func countoff() {
	for count := 1; count <= 2000; count++ {
		go func() {
			res, _ := cli.ExecutePlaybook(&request.TriggerPlaybook{
				AppToken:   "dXpuYTVjNzZpYWNuaThx",
				AppName:    "first_example",
				PlaybookId: 138,
				Data:       1000,
				Sync:       true,
			})

			b, err := json.Marshal(res)
			if err != nil {
				log.Println("err", err)
				return
			}
			log.Println(string(b))
		}()

	}

	select {}
}

type Payload struct {
	Info        string `json:"info"`
	TouchBizIds string `json:"touchBizIds"`
}

func countonline() {
	var cli = triggersdk.NewHttpClient("http://flow-service.xxx.com")

	for count := 1; count <= 10000; count++ {
		go func() {
			res, _ := cli.ExecutePlaybook(&request.TriggerPlaybook{
				AppToken:   "MXJ3ZHpobm5yanFueTh2",
				AppName:    "example",
				PlaybookId: 6,
				Data:       1000,
				Sync:       true,
			})

			b, err := json.Marshal(res)
			if err != nil {
				log.Println("err", err)
				return
			}
			log.Println(string(b))
		}()
	}

	select {}
}

func main() {
	//go loop()
	//for i := 0; i <= 1000; i++ {
	//	go count()
	//}
	//loop()
	countoff()
	//countonline()
	select {}
}
