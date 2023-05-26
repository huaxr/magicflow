// Author: XinRui Hua
// Time:   2022/2/8 下午3:45
// Git:    huaxr

package consumer

import (
	"context"
	"errors"
	"log"

	"strings"

	"github.com/huaxr/magicflow/component/dispatch"
	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/triggersdk"
	"github.com/huaxr/magicflow/pkg/workersdk/flow"
)

type WorkerConfig struct {
	Http        *triggersdk.HttpClient
	HandlerFunc Handler
	Concurrent  int
	MaxInFlight int
	Token       string
}

type Handler func(ctx flow.Context, msg *core.Msg) (interface{}, error)

func RegisterSlave(cx context.Context, namespace string, conf *WorkerConfig, mq dispatch.MQ) (chan *core.Event, error) {
	switch mq {
	case dispatch.Kafka:
		panic("err")

	case dispatch.NSQ:
		config := new(nsqconfig)
		config.nsqAddr = nil
		// get nsq
		d, err := conf.Http.GetLookUps()
		if err != nil {
			log.Printf("err when request lookups")
			return nil, err
		}
		if len(d) == 0 {
			return nil, errors.New("nsqconfig lookup not set")
		}
		// for every participator, lookups is the fundamental service discovery
		// for real brokers, nsq sdk using this functionality to start up consumer
		// thread and Handler is a exposure to receive data from connections.
		lookupAddrs := strings.Split(d, ",")
		for _, i := range lookupAddrs {
			config.lookupAddr = append(config.lookupAddr, strings.Trim(i, " "))
		}
		config.concurrent = conf.Concurrent
		consumer := new(nsqConsumer)
		consumer.w = &nsqworker{
			namespace:   namespace,
			token:       conf.Token,
			handlerFunc: conf.HandlerFunc,
			eventChan:   make(chan *core.Event, chanSize),
			maxinflight: conf.MaxInFlight,
		}
		consumer.http = conf.Http
		consumer.config = config
		consumer.ctx = cx
		consumer.topic = core.GetWorkerTopic(namespace)
		//slave.validate = validator.New()
		consumer.InitConsumer()
		consumer.Consume(cx)

		return consumer.w.eventChan, nil
	case dispatch.MOCK:
	default:

	}
	panic("")
}
