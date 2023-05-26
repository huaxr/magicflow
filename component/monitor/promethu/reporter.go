// Author: XinRui Hua
// Time:   2022/4/21 下午2:47
// Git:    huaxr

package promethu

import "context"

type Reporter int

const (
	Push Reporter = iota
	Pull
)

func LaunchProm(ctx context.Context, tp Reporter) {
	switch tp {
	case Pull:
		PromPull(ctx)
	case Push:
		PromPush(ctx)
	default:
		panic("not implement")
	}
}
