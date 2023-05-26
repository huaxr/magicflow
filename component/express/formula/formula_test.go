// Author: huaxr
// Time:   2021/7/5 上午11:19
// Git:    huaxr

package formula

import (
	"fmt"
	"github.com/tidwall/gjson"
	"regexp"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	const j = `{"name":{"first":"hua","last":"rui"},"age":1}`
	value := gjson.Get(j, "name.last")
	println(value.String())
}

// len($.ADD)==10 -> len([i])==10
func TestRe(t *testing.T) {
	text := "$.A.A==1 and len($.ADD)==10 and in($$$.EEE.field, $$.XXX.s)"
	reg := regexp.MustCompile(`(\$+\.[\.\w]+)`)
	res := reg.FindAllString(text, -1)

	for index, i := range res {
		text = strings.Replace(text, i, fmt.Sprintf("[param%v]", index), -1)
	}

	t.Logf("%+v", text)
}

func TestFormula(t *testing.T) {
	e := FormulaExpress{
		Exps: []*Expression{{
			Condition: "equals([param0], 'hello world')",
			Params:    []interface{}{"'hello world'"},
		}},
		Token: Default,
	}
	b, err := e.EExecute()
	t.Logf("%+v %v", b, err)

	e = FormulaExpress{
		Exps: []*Expression{{
			Condition: `equals([param0], 'hello world')`,
			Params:    []interface{}{"hello world"},
		}},
		Token: Default,
	}
	b, err = e.EExecute()
	t.Logf("%+v %v", b, err)

	e = FormulaExpress{
		Exps: []*Expression{{
			Condition: "equals([param0], [param1])",
			Params:    []interface{}{"x x", "x x"},
		}},
		Token: Default,
	}
	b, err = e.EExecute()
	t.Logf("%+v %v", b, err)

	e = FormulaExpress{
		Exps: []*Expression{{
			Condition: "equals(xxx, [param0])",
			Params:    []interface{}{"xxx"},
		}},
		Token: Default,
	}
	b, err = e.EExecute()
	t.Logf("%+v %v", b, err)

}

func TestFormula2(t *testing.T) {
	e := FormulaExpress{
		Exps: []*Expression{{
			Condition: `equals(getkey([param0], aaaa), 111)`,
			Params:    []interface{}{map[string]interface{}{"aaaa": 111}},
		}},
		Token: Default,
	}
	b, err := e.EExecute()
	t.Logf("%+v %v", b, err)
	select {}
}
