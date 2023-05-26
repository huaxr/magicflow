/**
 * @Author: huaxr
 * @Description:
 * @File: dcc_test
 * @Version: 1.0.0
 * @Date: 2021/7/28 下午4:55
 */

package dcc

import (
	"context"
	"encoding/json"
	"github.com/huaxr/magicflow/component/consensus"
	"github.com/huaxr/magicflow/core"
	"testing"
)

func TestDccPutWatch(t *testing.T) {
	pbId := 5
	ssId := 75
	globalDcc = new(dcc)
	globalDcc.con = consensus.NewConsensus(context.Background(), consensus.ETCD)

	err := globalDcc.DccPutPb(pbId, ssId)
	if err != nil {
		t.Log(err)
	}
}

func TestUpdateSnapshot(t *testing.T) {
	pb := core.NewPlayBookFromDb(1)
	res, err := json.Marshal(pb)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(string(res))
}
