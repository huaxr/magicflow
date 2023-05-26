// Author: huaxr
// Time:   2022/1/5 下午6:18
// Git:    huaxr

package core

import (
	"testing"

	"github.com/huaxr/magicflow/component/express"
)

func TestEdge(t *testing.T) {
	e := &edge{
		NextNode: "x",
		Express:  "$$. == 'A'",
		Priority: 0,
	}
	engine, _ := e.parseEngine("A", nil, express.Valuate)
	res, err := engine.EExecute()
	t.Log(res, err)

	e.Express = "$.A.B == 'A'"
	engine, _ = e.parseEngine("A", map[string]interface{}{"A": map[string]interface{}{"B": "A"}}, express.Valuate)
	res, err = engine.EExecute()
	t.Log(res, err)

	e.Express = "$.A.B == $$.B"
	engine, _ = e.parseEngine(map[string]interface{}{"B": "A"}, map[string]interface{}{"A": map[string]interface{}{"B": "A"}}, express.Valuate)
	res, err = engine.EExecute()
	t.Log(res, err)

	e.Express = "$.A.B == $$."
	engine, _ = e.parseEngine("A", map[string]interface{}{"A": map[string]interface{}{"B": "A"}}, express.Valuate)
	res, err = engine.EExecute()
	t.Log(res, err)

	e.Express = "equals($.A.B, 'A')"
	engine, _ = e.parseEngine("A", map[string]interface{}{"A": map[string]interface{}{"B": "A"}}, express.Valuate)
	res, err = engine.EExecute()
	t.Log(res, err)

	e.Express = "equals($.A.B, A)"
	engine, _ = e.parseEngine("A", map[string]interface{}{"A": map[string]interface{}{"B": "A"}}, express.Valuate)
	res, err = engine.EExecute()
	t.Log(res, err)

	e.Express = "xxx$$%%s "
	engine, _ = e.parseEngine("A", map[string]interface{}{"A": map[string]interface{}{"B": "A"}}, express.Valuate)
	res, err = engine.EExecute()
	t.Log(res, err)

	e.Express = "$$.x==1&&   $$.         ==2"
	engine, _ = e.parseEngine(map[string]interface{}{"x": 1, "yuuuu$": 2}, map[string]interface{}{"A": map[string]interface{}{"B": "A"}}, express.Valuate)
	res, err = engine.EExecute()
	t.Log(res, err)
}
