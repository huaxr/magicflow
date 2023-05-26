// Author: XinRui Hua
// Time:   2022/3/22 下午3:48
// Git:    huaxr

package client

import (
	"math/rand"
	"testing"
)

func TestServices(t *testing.T) {
	rpcServices.register("4")

	t.Log(rpcServices)

	rpcServices.offline("3")
	t.Log(rpcServices)

	rpcServices.offline("5")
	t.Log(rpcServices)
}

func TestRandom(t *testing.T) {

}
