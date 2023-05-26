// Author: huaxr
// Time:   2021/12/27 上午11:29
// Git:    huaxr

package valuate

import "testing"

func TestMains(t *testing.T) {
	engine := new(ValuateExpress)
	engine.Condition = "(equals(1,1) || gte(4, param)) && equals(1,1)"
	engine.Params = map[string]interface{}{"param": 9}
	res, err := engine.EExecute()
	t.Log(res, err)
}

func TestMains2(t *testing.T) {
	engine := new(ValuateExpress)
	engine.Condition = "1.0 << 1 == 1"
	engine.Params = nil
	res, err := engine.EExecute()
	t.Log(res, err)
}
