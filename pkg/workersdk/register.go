// Author: huaxr
// Time:   2021/6/25 下午3:47
// Git:    huaxr

package workersdk

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/huaxr/magicflow/component/dispatch"
	"github.com/huaxr/magicflow/component/dispatch/consumer"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/triggersdk"
	"github.com/huaxr/magicflow/pkg/workersdk/flow"
	"github.com/huaxr/magicflow/pkg/workersdk/registry"
	"github.com/spf13/cast"
)

// labor registe to the topic and knit one row, purl one row.
// doing the dispatcher message task.
// provider spectator for the Server's internal event.
func registerWorker(w *worker, mq dispatch.MQ) (chan *core.Event, error) {
	cli := triggersdk.NewHttpClient(w.server)
	err := w.validate(cli)
	if err != nil {
		log.Println("registerWorker.err", err)
		os.Exit(-1)
	}

	eveChan, err := consumer.RegisterSlave(w.ctx, w.app, &consumer.WorkerConfig{
		Http: cli,

		HandlerFunc: func(ctx flow.Context, msg *core.Msg) (interface{}, error) {
			// load configer here
			configer := msg.GetConfiguration()
			for _, comm := range configer.BeforeExecute {
				typ := cast.ToString(comm.Type)
				com, ok := registry.GetCommand(typ)
				if !ok {
					log.Printf("command not found:%v", typ)
					continue
				}
				com.Execute(ctx, comm.Paramters)
			}

			taskName := msg.GetTaskName()
			task, ok := registry.GetTask(taskName)
			if !ok {
				return nil, errors.New("task not exist:" + taskName)
			}

			var retry = configer.Retry

		gogo:
			var res *output
			var callback ofuture
			if configer.Timeout > 0 {
				// task.Do can be simultaneous with the handler commands
				// it can be a time lapse opration, wrap a future than
				// handling the callback at the end.
				callback = dofuturet(ctx, msg.GetInput(), time.Duration(configer.Timeout)*time.Second, task.Do)
			} else {
				callback = dofuture(ctx, msg.GetInput(), task.Do)
			}

			res = callback()
			// sweep last stride after callback executed successfully.
			// it may be err caused by slave exception or timeout.
			if res.err == nil {
				for _, comm := range configer.AfterExecute {
					typ := cast.ToString(comm.Type)
					com, ok := registry.GetCommand(typ)
					if !ok {
						log.Printf("command not found:%v", typ)
						continue
					}
					com.Execute(ctx, comm.Paramters)
				}
			} else {
				// when err not nil, and retry bigger than 0, just recall process
				if retry >= 1 {
					log.Printf("taskName:%v time:%v", taskName, retry)
					retry--
					goto gogo
				}
			}
			return res.output, res.err
		},
		Concurrent:  w.concurrent,
		Token:       w.token,
		MaxInFlight: w.batch,
	}, mq)

	if err != nil {
		log.Printf("registerWorker err %v", err)
		return nil, err
	}

	log.Printf("register success, welcome, %v", w.app)

	defer registry.LoadHB()
	return eveChan, nil
}
