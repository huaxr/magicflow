// Author: huaxr
// Time:   2022/1/4 下午2:57
// Git:    huaxr

package publisher

import (
	"context"
	"io"
	"io/ioutil"
)

type mockPublisher struct {
	topicchan map[string]chan []byte
}

func (p *mockPublisher) Publish(topic, broker string, in io.Reader, infors ...interface{}) (err error) {
	body, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	p.topicchan[topic] <- body
	return nil
}

func (p *mockPublisher) AvailableBrokersCount() int32 {
	return 1
}

func (p *mockPublisher) Heartbeat()   {}
func (p *mockPublisher) Close() error { return nil }

func InitMOCKPublisher(ctx context.Context) *mockPublisher {
	return nil
}
