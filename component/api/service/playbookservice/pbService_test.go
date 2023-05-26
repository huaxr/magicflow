/**
 * @Author: huaxr
 * @Description:
 * @File: pbService_test
 * @Version: 1.0.0
 * @Date: 2021/9/1 下午5:23
 */

package playbookservice

import (
	"context"
	"testing"

	"github.com/huaxr/magicflow/core"
)

func TestCreatePlayBook(t *testing.T) {
	req := CreatePlayBookReq{
		App:         "test_playbook_api",
		Name:        "tmp",
		Description: "测试剧本api",
	}
	res, err := CreateEmptyPlayBook(context.Background(), &req)
	if err != nil {
		t.Log(err)
	}
	t.Log(res)
	pbs := core.G()
	pb := pbs.GetPlaybook(res.PlayBookId)
	t.Log(pb)
}
