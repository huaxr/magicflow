// Author: huaxr
// Time:   2021/8/6 下午2:12
// Git:    huaxr

package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"unsafe"

	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/request"
	"github.com/huaxr/magicflow/pkg/triggersdk"
	"github.com/huaxr/magicflow/pkg/workersdk/registry"
	"github.com/nsqio/go-nsq"
)

var (
	nsqLogger = log.New(ioutil.Discard, "", log.LstdFlags)
)

const (
	chanSize = 1 << 16
)

type nsqworker struct {
	token       string
	namespace   string
	handlerFunc Handler
	maxinflight int
	slave       *unsafe.Pointer
	// nsqworker register and listen specific event.
	eventChan chan *core.Event
}

type nsqconfig struct {
	nsqAddr    []string
	lookupAddr []string
	concurrent int
}

type nsqConsumer struct {
	ctx    context.Context
	topic  string
	config *nsqconfig
	w      *nsqworker
	http   *triggersdk.HttpClient
	//validate    *validator.validate
}

func (h *nsqConsumer) handlerMessage(body []byte) (err error) {
	var msg core.Msg
	err = json.Unmarshal(body, &msg)
	if err != nil {
		//handlerBadMsg(m.Body, "Unmarshal json failed.")
		log.Printf("slave HandleMessage err %v", err.Error())
		return
	}

	//t1 := time.Now().UnixNano()
	// output could be anything or limits?
	output, err := h.w.handlerFunc(msg.Env, &msg)
	//t2 := time.Now().UnixNano()

	// the nsqworker can set output nil cause there a situations
	// that the node need some mutual interaction with a person
	// or a unexpected option.

	// report exceptions to server.
	if err != nil {
		_, err = h.http.ReportException(&request.WorkerExceptionReq{
			ServiceAddr: msg.ServiceAddr,
			Key:         msg.Key,
			Signature:   msg.Signature,
			Exception:   err.Error(),
			Token:       h.w.token,
		})
		return
	}

	// allowed
	if output == nil {
		// nsqworker response api need bingding the  output field.
		output = ""
	}

	// Determine output whether serializable
	_, err = json.Marshal(output)
	// report exceptions to server.
	if err != nil {
		_, err = h.http.ReportException(&request.WorkerExceptionReq{
			ServiceAddr: msg.ServiceAddr,
			Key:         msg.Key,
			Signature:   msg.Signature,
			Exception:   "output is not serializable",
			Token:       h.w.token,
		})
		return
	}

	// Handler nsqworker trigger.
	// notice: server no longer be a slave again.
	_, err = h.http.ReportTask(&request.WorkerResponseReq{
		ServiceAddr: msg.ServiceAddr,
		Key:         msg.Key,
		Signature:   msg.Signature,
		Output:      output,
		Token:       h.w.token,
		HeartBeat:   registry.HB,
	})
	return
}

func (h *nsqConsumer) handlerEvent(body []byte) (err error) {
	var event core.Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		//handlerBadMsg(m.Body, "Unmarshal json failed.")
		log.Printf("slave HandleMessage err %v", err.Error())
		return
	}
	if len(h.w.eventChan) < cap(h.w.eventChan) {
		h.w.eventChan <- &event
		return
	}
	log.Printf("handlerEvent eventChan buffer full")
	return
}

// HandleMessage is implement of nsq.Consume.Handler.
// which handle the nsq.message and impetus the in-flight loop works.
// when return error != nil, OnRequeue will be triggered. else OnFinish
// will tell the nsqd while !m.DisableAutoResponse()
func (h *nsqConsumer) HandleMessage(m *nsq.Message) (err error) {
	//m.DisableAutoResponse()
	body, typ := m.Body[:len(m.Body)-1], m.Body[len(m.Body)-1:][0]
	switch typ {
	case core.MESSAGE:
		go func() {
			err = h.handlerMessage(body)
		}()
	case core.EVENT:
		err = h.handlerEvent(body)
	default:
		log.Printf("recognition fail, %v", typ)
	}
	return
}

// InitConsumer nsq worker as a slave and send response to server. server become a distribute center.
func (h *nsqConsumer) InitConsumer() {
	config := nsq.NewConfig()
	// the mq cluster configurated the auth-url with tmp-env address link.
	// the /auth api will never reached in local environment.
	config.AuthSecret = fmt.Sprintf("%v?%v", h.w.token, h.w.namespace)

	if h.w.maxinflight != 0 {
		config.MaxInFlight = h.w.maxinflight
	}
	config.Snappy = true
	consumer, err := nsq.NewConsumer(h.topic, "nil", config)
	if err != nil {
		panic(err)
	}
	consumer.SetLogger(nsqLogger, nsq.LogLevelDebug)
	// slave -> &slave will panic.
	// the parameter of unsafe.Pointer is pointer already.
	h.setConsumer(unsafe.Pointer(consumer))
}

func (h *nsqConsumer) connectDirectly() {
	err := h.getConsumer().ConnectToNSQDs(h.config.nsqAddr)
	if err != nil {
		panic(err)
	}
	stats := h.getConsumer().Stats()
	if stats.Connections == 0 {
		panic("stats report 0 connections (should be > 0)")
	}
}

func (h *nsqConsumer) connectLookUp() {
	c := h.getConsumer()
	if err := c.ConnectToNSQLookupds(h.config.lookupAddr); err != nil {
		log.Printf("connectLookUp err %v", err)
		panic(err)
	}
}

func (h *nsqConsumer) Consume(ctx context.Context) {
	h.getConsumer().AddConcurrentHandlers(h, h.config.concurrent)
	h.connectLookUp()
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("stop nsq slave&pool.")
				h.Close()
				return
			}
		}
	}()
	log.Printf("start %d nsq consumers", h.config.concurrent)
}

func (h *nsqConsumer) Close() {
	h.getConsumer().Stop()
}

func (h *nsqConsumer) getConsumer() *nsq.Consumer {
	return (*nsq.Consumer)(*h.w.slave)
}

func (h *nsqConsumer) setConsumer(pointer unsafe.Pointer) {
	h.w.slave = &pointer
}
