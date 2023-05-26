// Author: huaxr
// Time:   2021/10/8 下午5:11
// Git:    huaxr

package core

import (
	"reflect"
	"strings"

	"github.com/huaxr/magicflow/component/express/parser"
)

type ttlTyp int

// referenced by playbook context node edge.
type ttl map[string]int

const (
	_ ttlTyp = iota

	// remember all the paths executes result in the context
	// without any strategy to handle context slim.
	remember
	// using cache by global ttl.
	// this strategy need one global ttl in pb and node ttl
	// for restore the context while exiting the node.
	// it's a crude method of pre - calculation, which not support
	// multi branches result merge.

	// combine message.Context.TTL with g-ttl. desperate now.
	// for the sake of context expand rapidly as we define a playbook with multiply
	// node, for each of them stored in the context will cause severe performance
	// degradation(spd), so we leads to a solution which calculating by likelihood estimation
	// the reference count of node by algorithm, till we get a moderate size.
	// the optimum slime of context set before initialized parse.
	// this is implement by quote. now replaced by dp.
	quote
	// dynamic planning arithmetic.
	// this is an optimal path selection strategy, which need
	// every branch remember delete what context key when entering
	// the next node. (actually compare each reference and do difference set)
	// it differs from quote strategy in that the input is specified
	// as a global node through dynamic planning, and then the
	// change path is told to delete and take some data to save
	/*
		        +--[equals(1,1)]---> B
			A --|
				+---[gte($.A.NUM, 10)]---> C -- > D($.C)
	*/
	// branch A->B told that A dose not remember any more
	//	- A ttl is {}
	// branch A->C told that B dose not used in the later
	// 	- A ttl is {"A":1}
	// branch C->D told C must record C's output
	//	- C ttl is {"C":1} so delete A from context.
	dp
)

type slimer struct {
	ttltyp ttlTyp
	ttl    ttl
}

func merge(tts ...ttl) ttl {
	var ttl = make(ttl)
	for _, tt := range tts {
		if len(tt) == 0 {
			continue
		}
		for k := range tt {
			ttl[k] += 1
		}
	}
	return ttl
}

// referer TestTTL
func ttlx(v interface{}, tt ttl) ttl {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Map:
		for _, vv := range v.(map[string]interface{}) {
			tt = ttlx(vv, tt)
		}

	case reflect.String:
		// panic: interface conversion: interface {} is json.Number, not string
		// give a assert here
		if _, ok := v.(string); ok {
			v = strings.Trim(v.(string), " ")
			res := reg.FindAllString(v.(string), -1) // ["$.xx.yy"]
			for _, s := range res {
				// should not ignore $$$.
				if strings.HasPrefix(s, parser.TriggerPrefix) {
					tt[parser.TriggerContextKey] = 1
					continue
				}
				// ignore $$.X && $._env
				if strings.HasPrefix(s, parser.EnvPrefix) {
					continue
				}
				if strings.HasPrefix(s, parser.ContextPrefix) {
					if x := strings.Split(s, parser.Point); len(x) > 1 {
						tt[x[1]] = 1
					}
				}
			}
		}
	case reflect.Array, reflect.Slice:
		for _, vv := range v.([]interface{}) {
			tt = ttlx(vv, tt)
		}
	}

	return tt
}
