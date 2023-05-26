// Author: XinRui Hua
// Time:   2022/1/27 下午3:05
// Git:    huaxr

package consumer

import (
	"context"
)

type mockConsumer struct {
}

func (h *mockConsumer) Consume(ctx context.Context) {

}

func (h *mockConsumer) InitConsumer() {

}
