// Author: huaxr
// Time:   2021/8/2 下午7:33
// Git:    huaxr

package dcc

import (
	"context"

	"github.com/huaxr/magicflow/component/consensus"
	"github.com/huaxr/magicflow/pkg/accutil"
)

func GetDcc() *dcc {
	if globalDcc != nil {
		return globalDcc
	}
	panic("globalDcc is nil")
}

func LaunchDcc(ctx context.Context, typ consensus.KVType) {
	globalDcc = new(dcc)
	globalDcc.con = consensus.NewConsensus(ctx, typ)
	globalDcc.ctx = ctx
	// start watching
	go accutil.Thread("watch-dcc", func() { globalDcc.dccWatch() })
}
