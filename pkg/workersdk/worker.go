// Author: huaxr
// Time:   2021/7/19 下午7:11
// Git:    huaxr

package workersdk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"time"

	"github.com/huaxr/magicflow/component/dispatch"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/request"
	"github.com/huaxr/magicflow/pkg/triggersdk"
	"github.com/huaxr/magicflow/pkg/workersdk/flow"
	"github.com/spf13/cast"
)

type (
	// asynchronous task executes
	ifuture func(ctx flow.Context, input interface{}) (output interface{}, err error)
	ofuture func() *output

	worker struct {
		ctx        context.Context
		server     string
		app        string
		token      string
		concurrent int
		// max in flight
		batch int
		// info 	0001
		// cliError 0010
		// srvError 0100
		// notify   1000
		eventFlag uint8
	}

	ConfigWorker struct {
		FlowServer  string
		App         string
		Token       string
		Concurrency int
		// adjust until reached a optimal plateau.
		Batch int
	}

	output struct {
		output interface{}
		err    error
	}
)

// NewWorker
func NewWorker(ctx context.Context, conf *ConfigWorker) *worker {
	w := new(worker)
	w.ctx = ctx
	w.server = conf.FlowServer
	w.app = conf.App
	// The token will guarantee you successfully get the secret key
	// and get NSQ's broker connection permission.
	w.token = conf.Token
	w.concurrent = conf.Concurrency
	w.eventFlag = 0xf
	w.batch = conf.Batch
	return w
}

func dofuturet(ctx flow.Context, input interface{}, dur time.Duration, f ifuture) ofuture {
	var (
		result interface{}
		err    error
		c      = make(chan struct{}, 0)
		// for interaction notification
		done = make(chan struct{}, 0)
	)

	go time.AfterFunc(dur, func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()

		select {
		case _, _ = <-done:
			return
		default:
			err = fmt.Errorf("time exceed:%v", dur)
		}
		c <- struct{}{}
	})

	go func() {
		defer func() {
			if err := recover(); err != nil {
				return
			}
		}()

		select {
		case <-done:
			return
		default:
			result, err = f(ctx, input)
		}
		c <- struct{}{}
	}()

	return func() *output {
		<-c
		close(c)
		close(done)
		return &output{
			output: result,
			err:    err,
		}
	}
}

func dofuture(ctx flow.Context, input interface{}, f ifuture) ofuture {
	var (
		result interface{}
		err    error
		c      = make(chan struct{}, 0)
	)

	go func() {
		var er error
		result, er = f(ctx, input)
		err = er
		c <- struct{}{}
	}()

	return func() *output {
		<-c
		close(c)
		return &output{
			output: result,
			err:    err,
		}
	}
}

func (w *worker) validate(cli *triggersdk.HttpClient) error {
	// impose restrictions on the namespace & token validation.
	if len(w.app) > 50 || len(w.app) < 3 {
		return fmt.Errorf("format app err:%v", "namespace's length range [3,50]")
	}

	var result request.HasAuthRes
	resp, err := cli.GetRestClient().R().
		SetHeader("Content-Type", "application/json").
		SetBody(&request.WorkerAuth{
			Namespace: w.app,
			Token:     w.token,
		}).
		Post(fmt.Sprintf("%s%s", cli.GetBasePath(), "/config/auth"))

	if err != nil {
		log.Printf("validate worker err %v", err)
		return err
	}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return err
	}

	if result.Data["status"] != "ok" {
		return fmt.Errorf("auth err: %v", "not authentication!")
	}

	w.token = cast.ToString(result.Data["secret"])

	if w.batch < cast.ToInt(result.Data["broker_size"]) {
		resize := cast.ToInt(result.Data["broker_size"]) + 1
		log.Printf("reset batch from %v to %v", w.batch, resize)
		w.batch = resize
	}

	if w.batch > 1e4 {
		log.Printf("batch size too big, you should judge your app eps, set a rational batch size to 10000")
		w.batch = 1e4
	}
	return nil
}

func (w *worker) Register() (chan *core.Event, error) {
	return registerWorker(w, dispatch.NSQ)
}

// before consumer
func (w *worker) BeforeFuncs() {
	panic("not implemnet")
}

func (w *worker) AfterFuncs() {
	panic("not implemnet")
}
