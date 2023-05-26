// Author: huaxr
// Time:   2022/1/26 上午11:41
// Git:    huaxr

package publisher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/huaxr/magicflow/component/ticker"

	"github.com/huaxr/magicflow/component/logx"

	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/huaxr/magicflow/pkg/accutil"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/Shopify/sarama"
)

var (
	kafkaProducer *kafkaPublisher
)

type kafkaPublisher struct {
	sync.RWMutex
	ctx context.Context
	// binding address with client
	availableAddress map[string]sarama.AsyncProducer
	// flag of initiation
	initialized int32
}

func InitKafkaPublisher(ctx context.Context) *kafkaPublisher {
	kafkaProducer = new(kafkaPublisher)
	kafkaProducer.initialized = 0
	kafkaProducer.ctx, _ = context.WithCancel(ctx)

	kafkaProducer.availableAddress = make(map[string]sarama.AsyncProducer)
	ticker.RegisterTick(kafkaProducer)
	return kafkaProducer
}

func (p *kafkaPublisher) Publish(topic, broker string, in io.Reader, infors ...interface{}) (err error) {
	body, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	// send message
	msg := &sarama.ProducerMessage{
		Topic: topic,
		// hash this key
		Key:   sarama.StringEncoder(infors[0].(string)),
		Value: sarama.ByteEncoder(body),
		// when config.Producer.Partitioner = sarama.NewManualPartitioner
		// Partition: cast.ToInt32(partition),
	}

	var count int32
RETRY:
	if p.initialized == 1 {
		if broker != "" {
			p.RLock()
			if pp, ok := p.availableAddress[broker]; ok {
				pp.Input() <- msg
			} else {
				logx.L().Errorf("broker %v lost", broker)
			}
			p.RUnlock()
		}
	} else {
		logx.L().Warnf("publisher not initialized, now ticker")
		p.ticker()
		count++
		if count >= p.AvailableBrokersCount() {
			logx.L().Errorf("publisher dead for %d times", count)
			return
		}
		goto RETRY
	}
	return
}

func (p *kafkaPublisher) AvailableBrokersCount() int32 {
	return int32(len(p.availableAddress))
}

func (p *kafkaPublisher) Heartbeat() {
	p.ticker()
	job := ticker.NewJob(p.ctx, p.Name(), p.Duration(), func() {
		p.Lock()
		p.ticker()
		p.Unlock()
	})
	ticker.GetManager().Register(job)
}

func (p *kafkaPublisher) Name() string           { return "kafka_heartbeat" }
func (p *kafkaPublisher) Duration() *time.Ticker { return time.NewTicker(5 * time.Second) }

func (p *kafkaPublisher) ticker() {
	availableAddress := strings.Split(confutil.GetConf().Queue.Kafka.Brokers, ",")
	//availableAddress := []string{"10.90.73.26:9092", "10.90.73.54:9092", "10.90.73.56:9092"}
	if len(p.availableAddress) == len(availableAddress) {
		return
	}

	for _, addr := range availableAddress {
		addr = strings.TrimSpace(addr)
		if _, ok := p.availableAddress[addr]; ok {
			continue
		}
		p.addProducer(addr, p.connect(addr))
	}

	if len(p.availableAddress) == 0 {
		accutil.LockSwap(&p.initialized, 1, 0)
		logx.L().Warnf("no available pool found. heartbeat trigger next ticker.")
	} else {
		accutil.LockSwap(&p.initialized, 0, 1)
	}
}

func (p *kafkaPublisher) addProducer(addr string, pro sarama.AsyncProducer) {
	if pro == nil {
		return
	}
	p.availableAddress[addr] = pro
	logx.L().Infof("start/recover %s kafka publisher", addr)

}

func (p *kafkaPublisher) connect(addr string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	// define manual
	//config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V0_11_0_2

	pro, err := sarama.NewAsyncProducer([]string{addr}, config)
	if err != nil {
		logx.L().Errorf("kafkaPublisher connect %v", err)
		return nil
	}

	go func() {
		for {
			select {
			case suc := <-pro.Successes():
				fmt.Printf("offset: %d,  timestamp: %s", suc.Offset, suc.Timestamp.String())
			case fail := <-pro.Errors():
				logx.L().Errorf("err %v, addr:%v", fail.Err.Error(), addr)
			}
		}
	}()
	return pro
}

func (p *kafkaPublisher) Close() error {
	p.Lock()
	for _, pr := range p.availableAddress {
		pr.AsyncClose()
	}
	p.availableAddress = map[string]sarama.AsyncProducer{}
	p.initialized = 0
	p.Unlock()

	return nil
}

func (p *kafkaPublisher) PublishWithRetry(topic, broker string, body []byte, try int32) error {
	if try == 0 {
		logx.L().Warnf("kafkaPublisher try excessed")
		return nil
	}

	err := p.Publish(topic, broker, bytes.NewBuffer(body))
	if err != nil {
		atomic.AddInt32(&try, -1)
		logx.L().Warnf("kafkaPublisher untilTryout err: %s, topic: %s, leftTime: %v", err.Error(), topic, try)
		return p.PublishWithRetry(topic, broker, body, try)
	}
	return err

}
