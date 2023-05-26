// Author: huaxr
// Time:   2021/8/12 上午11:56
// Git:    huaxr

package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/huaxr/magicflow/component/express"
	"github.com/huaxr/magicflow/component/express/formula"
	"github.com/huaxr/magicflow/component/express/parser"
	"github.com/huaxr/magicflow/component/express/valuate"
)

// Branches sorted by Priority
// Each node has one or more conditional edges, parallel edges, or restricted edges
type (
	branches   []*edge
	branchType string
)

const (
	// BranchChoice are mutually exclusive and unique
	BranchChoice   branchType = "choice"
	BranchParallel branchType = "parallel"
)

// branches belongs to one node.
// getTtl return references to all conditional edges of the node cites.
// when calc ttl, it should be considered.
func (s branches) getTtl() ttl {
	if len(s) == 0 {
		return ttl{}
	}

	var tt = make(ttl)
	for _, e := range s {
		res := reg.FindAllString(e.Express, -1)
		for _, s := range res {
			// should not ignore $$$.
			if strings.HasPrefix(s, parser.TriggerPrefix) {
				tt[parser.TriggerContextKey] = 1
				continue
			}
			// ignore $$.X
			if !strings.HasPrefix(s, parser.LastPrefix) {
				continue
			}
			if x := strings.Split(s, parser.Point); len(x) > 1 {
				tt[x[1]] = 1
			}
		}
	}
	return tt
}

func (s branches) Len() int           { return len(s) }
func (s branches) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s branches) Less(i, j int) bool { return s[i].Priority < s[j].Priority }

// Edge is the relations between the nodes, which contains decisive
// bool expression to mobile message orientation schedule.
// Condition field is the javascript expression.
// @2021.06.24 for the consideration of parallel executions without
// make things complicated, edge should support must-execute-once
// to enable other nodes execute directly in order to relief dag architecture.
type edge struct {
	// depict the usage and give a vivid to the unintelligible express
	// Description string `json:"description"`

	// NextNode implicit the edge end point. precede node disabled meanwhile.
	// recipient node binding with edge.
	NextNode string `json:"next_node"`
	// refer to json-util.GetKey to analyze $. expression.
	Express string `json:"express"`
	// Priority defined if-else code block order.
	Priority int `json:"priority"`
}

// pre nodes tells node must wait someone done.
type depends []string

func (r depends) size() int {
	return len(r)
}

// $.store.book[?(@.price <.expensive)].price
// regexp reference https://blog.csdn.net/weixin_39622138/article/details/113075233
var reg = regexp.MustCompile(`\${1,3}[\w|\.]+(\[(.*?)\])?[\w|\.]*`)

// Support dsl.And, dsl.Or
func (e *edge) parseFormulaExpress(last interface{}, mapper map[string]interface{}) (express.EngineImpl, error) {
	engine := new(formula.FormulaExpress)

	if strings.Contains(e.Express, string(formula.And)) {
		engine.Token = formula.And
	} else if strings.Contains(e.Express, string(formula.Or)) {
		engine.Token = formula.Or
	} else {
		engine.Token = formula.Default
	}

	conditions := strings.Split(e.Express, string(engine.Token))
	for _, condition := range conditions {
		var params = make([]interface{}, 0)

		res := reg.FindAllString(condition, -1)
		for index, i := range res {
			condition = strings.Replace(condition, i, fmt.Sprintf("[param%v]", index), -1)
			value, err := parser.GetKey(i, last, mapper)
			if err != nil {
				return nil, err
			}
			params = append(params, value)
		}
		engine.Exps = append(engine.Exps, &formula.Expression{Condition: condition, Params: params})
	}
	return engine, nil
}

func (e *edge) parseValuateExpress(last interface{}, mapper map[string]interface{}) (express.EngineImpl, error) {
	condition := e.Express
	var (
		params = make(map[string]interface{}, 0)
		err    error
	)
	res := reg.FindAllString(condition, -1)
	for index, dollarExp := range res {
		key := fmt.Sprintf("param%v", index)
		condition = strings.Replace(condition, dollarExp, key, -1)
		params[key], err = parser.GetKey(dollarExp, last, mapper)
		if err != nil {
			return nil, err
		}
	}
	return &valuate.ValuateExpress{
		Condition: condition,
		Params:    params,
	}, nil
}

func (e *edge) parseEngine(last interface{}, mapper map[string]interface{}, method express.Engine) (express.EngineImpl, error) {
	switch method {
	case express.Formula:
		return e.parseFormulaExpress(last, mapper)
	case express.Valuate:
		return e.parseValuateExpress(last, mapper)
	case express.Js:

	}
	return nil, fmt.Errorf("method %v not implement", method)
}
