// Author: huaxr
// Time:   2021/6/9 下午3:50
// Git:    huaxr

package publisher

import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/huaxr/magicflow/component/ticker"

	"github.com/huaxr/magicflow/component/logx"

	"io/ioutil"
	"log"
	"strings"
	"time"
	"unsafe"

	"github.com/huaxr/magicflow/pkg/accutil"
	"github.com/huaxr/magicflow/pkg/confutil"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/cast"
)

var (
	nsqProducer *nsqPublisher
	nsqLogger   = log.New(ioutil.Discard, "", log.LstdFlags)
)

// nsqPublisher provides function to dispense message to specified topic
// it's only works on server dispatch engine.
type nsqPublisher struct {
	sync.RWMutex
	ctx context.Context
	// binding address with client
	availableAddress map[string]unsafe.Pointer
	// flag of initiation
	initialized int32
}

func InitNSQPublisher(ctx context.Context) *nsqPublisher {
	nsqProducer = new(nsqPublisher)
	nsqProducer.initialized = 0
	nsqProducer.ctx, _ = context.WithCancel(ctx)

	nsqProducer.availableAddress = make(map[string]unsafe.Pointer)
	ticker.RegisterTick(nsqProducer)
	return nsqProducer
}

func (p *nsqPublisher) AvailableBrokersCount() int32 {
	return int32(len(p.availableAddress))
}

func (p *nsqPublisher) Heartbeat() {
	p.ticker()
	job := ticker.NewJob(p.ctx, p.Name(), p.Duration(), p.ticker)
	ticker.GetManager().Register(job)
}

func (p *nsqPublisher) Name() string { return "nsq_heartbeat" }

func (p *nsqPublisher) Duration() *time.Ticker {
	t := cast.ToInt(strings.TrimRight(confutil.GetConf().Configuration.BrokerHeartbeatInterval, "s"))
	return time.NewTicker(time.Duration(t) * time.Second)
}

// DELETE http://10.90.72.135:4171/api/nodes/10.90.72.172%3A4151
// BODY: {"topic": "xesFlow_first_example"}
func (p *nsqPublisher) delete(broker string, topic string) {

}

// for the sake of supervising mq cluster healthy and no offence to the
// nsq producer work.
// ticker when broker exit, tryUntil will raise err and republish message
// to the alive broker, ticker will monitor the pool and try ping the
// address provide in the config files.
func (p *nsqPublisher) ticker() {
	availableAddress := strings.Split(confutil.GetConf().Queue.Nsq.Brokers, ",")
	if len(p.availableAddress) == len(availableAddress) {
		// implicit some pool client broken.
		// if all the nsq pool client alive, just bypass.
		return
	}

	for _, addr := range availableAddress {
		addr = strings.TrimSpace(addr)
		if _, ok := p.availableAddress[addr]; ok {
			// only dial the died connection.
			continue
		}

		p.addProducer(p.connect(addr))
	}

	if len(p.availableAddress) == 0 {
		accutil.LockSwap(&p.initialized, 1, 0)
		logx.L().Warnf("no available pool found. heartbeat trigger next ticker.")
	} else {
		accutil.LockSwap(&p.initialized, 0, 1)
	}
}

func (p *nsqPublisher) Close() error {
	//Close(p.internalchan)
	p.Lock()
	defer p.Unlock()
	if len(p.availableAddress) == 0 {
		return fmt.Errorf("no avaliable address found")
	}
	for _, pr := range p.availableAddress {
		pr2 := (*nsq.Producer)(pr)
		if pr2.Ping() != nil {
			pr2.Stop()
		}
	}

	p.availableAddress = map[string]unsafe.Pointer{}
	p.initialized = 0
	return nil
}

func (p *nsqPublisher) addProducer(pro *nsq.Producer) {
	if pro == nil {
		return
	}
	p.Lock()
	p.availableAddress[pro.String()] = unsafe.Pointer(pro)
	p.Unlock()
	logx.L().Infof("start/recover %s nsq publisher", pro.String())
}

func (p *nsqPublisher) connect(addr string) *nsq.Producer {
	config := nsq.NewConfig()
	config.AuthSecret = confutil.GetConf().Queue.Nsq.Secret
	pro, err := nsq.NewProducer(addr, config)
	if err != nil {
		logx.L().Errorf("nsqPublisher err %v", err)
		return nil
	}
	// pro.ping will has nsq log
	pro.SetLogger(nsqLogger, nsq.LogLevelInfo)
	if err = pro.Ping(); err != nil {
		logx.L().Warnf("addr:%v ping err: %v", pro.String(), err)
		return nil
	}
	return pro
}

// dispatch AllowedBrokers dispatch message with the specified brokers
// with randomise & roundRobin strategy. the brokers should reserve quota by the app.
// familiar with knit one row, purl one row.
func (p *nsqPublisher) Publish(topic, broker string, in io.Reader, infos ...interface{}) (err error) {
	body, err := ioutil.ReadAll(in)
	if err != nil {
		logx.L().Errorf("nsq publish err:%v", err)
		return err
	}
	if len(body) == 0 || len(body) > 1048576 {
		return fmt.Errorf("dispatch err size: %v", len(body))
	}
	var count int32
RETRY:
	if p.initialized == 1 {
		p.RLock()
		pp, ok := p.availableAddress[broker]
		p.RUnlock()

		if ok {
			err = (*nsq.Producer)(pp).Publish(topic, body)
			if err != nil {
				count++
				if count >= p.AvailableBrokersCount() {
					logx.L().Errorf("publisher dead for %d times", count)
					return fmt.Errorf("nsq publish err:%v", err)
				}
				goto RETRY
			}
		} else {
			logx.L().Errorf("broker %v lost", broker)
		}
	} else {
		logx.L().Warnf("publisher not initialized, now ticker")
		p.ticker()
		goto RETRY
	}
	return
}

func (p *nsqPublisher) Export(string) []byte {
	var available = make([]string, 0)
	for key, _ := range p.availableAddress {
		available = append(available, key)
	}
	b, _ := json.Marshal(available)
	return b
}
