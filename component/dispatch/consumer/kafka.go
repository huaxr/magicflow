// Author: XinRui Hua
// Time:   2022/1/27 下午2:44
// Git:    huaxr

package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/huaxr/magicflow/core"
	"github.com/huaxr/magicflow/pkg/request"
	"github.com/huaxr/magicflow/pkg/triggersdk"
	"github.com/huaxr/magicflow/pkg/workersdk/registry"
	"github.com/Shopify/sarama"
)

type kafkaconfig struct {
	hosts []string
}

type kafkaworker struct {
	token       string
	namespace   string
	handlerFunc Handler
	maxinflight int
	slave       sarama.PartitionConsumer
	// nsqworker register and listen specific event.
	eventChan chan *core.Event
}

type kafkaConsumer struct {
	ctx    context.Context
	topic  string
	config *kafkaconfig
	w      *kafkaworker
	http   *triggersdk.HttpClient
}

// InitConsumer nsqworker as a slave and send response to server. server become a distribute center.
func (h *kafkaConsumer) InitConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V0_11_0_2

	// consumer
	consumer, err := sarama.NewConsumer(h.config.hosts, config)
	if err != nil {
		log.Printf("consumer_test create consumer error %s\n", err.Error())
		panic(err)
	}

	partitionconsumer, err := consumer.ConsumePartition(h.topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("try create partition_consumer error %s\n", err.Error())
		panic(err)
	}

	h.setConsumer(partitionconsumer)
}

func (h *kafkaConsumer) setConsumer(pointer sarama.PartitionConsumer) {
	h.w.slave = pointer
}

func (h *kafkaConsumer) Consume(ctx context.Context) {
	for {
		select {
		case msgx := <-h.w.slave.Messages():
			//fmt.Printf("msg offset: %d, partition: %d, timestamp: %s, value: %s\n",
			//	msgx.Offset, msgx.Partition, msgx.Timestamp.String(), string(msgx.Value))
			var msg core.message
			err := json.Unmarshal(msgx.Value, &msg)
			if err != nil {
				//handlerBadMsg(m.Body, "Unmarshal json failed.")
				log.Printf("slave HandleMessage err %v", err.Error())
				return
			}
			output, err := h.w.handlerFunc(msg.Context.GetEnv(), &msg)
			if err != nil {
				_, err = h.http.ReportException(&request.WorkerExceptionReq{
					Message:   &msg,
					Exception: err.Error(),
					Token:     h.w.token,
				})
				return
			}
			if output == nil {
				output = ""
			}

			_, err = json.Marshal(output)
			if err != nil {
				_, err = h.http.ReportException(&request.WorkerExceptionReq{
					Message:   &msg,
					Exception: "output is not serializable",
					Token:     h.w.token,
				})
				return
			}
			_, err = h.http.ReportTask(&request.WorkerResponseReq{
				Message:   &msg,
				Output:    output,
				Token:     h.w.token,
				HeartBeat: registry.HB,
			})
			return

		case err := <-h.w.slave.Errors():
			fmt.Printf("err :%s\n", err.Error())
		}
	}
}
