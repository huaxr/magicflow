// Author: huaxr
// Time:   2021/8/5 上午9:56
// Git:    huaxr

package dispatch

import (
	"context"
	"io"

	"github.com/huaxr/magicflow/component/dispatch/publisher"
)

type MQ string

const (
	NSQ   MQ = "nsq"
	Kafka MQ = "kafka"
	MOCK  MQ = "mock"
)

// multilateral support of Publisher, including nsq, kafka, tmp.
type Publisher interface {
	io.Closer
	// publish a message or event.
	Publish(topic, broker string, in io.Reader, extra ...interface{}) (err error)
	AvailableBrokersCount() int32
	// watch the available broker connection.
	Heartbeat()
}

type Consumer interface {
	Consume(ctx context.Context)
	InitConsumer()
}

func InitMQ(ctx context.Context, mt MQ) Publisher {
	switch mt {
	case NSQ:
		return publisher.InitNSQPublisher(ctx)

	case Kafka:
		return publisher.InitKafkaPublisher(ctx)

	case MOCK:
		return publisher.InitMOCKPublisher(ctx)

	}
	panic("not implement yet")
}
