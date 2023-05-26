// Author: huaxr
// Time:   2021/7/2 上午11:32
// Git:    huaxr

package parser

import (
	"encoding/json"
	"fmt"
	"github.com/oliveagle/jsonpath"
	"testing"
)

func TestGetKey(t *testing.T) {
	var res interface{}
	var err error
	res, _ = GetKey("  $.AAA.number", "", map[string]interface{}{"AAA": map[string]interface{}{"number": 1}})
	t.Logf("%v", res)

	res, _ = GetKey("$.AAA", "", map[string]interface{}{"AAA": map[string]interface{}{"number": 1}})
	t.Logf("%v", res)

	res, _ = GetKey("$.", "", map[string]interface{}{"AAA": map[string]interface{}{"number": 1}})
	t.Logf("%v", res)

	res, _ = GetKey("$.AAA.numbers[2]", "", map[string]interface{}{"AAA": map[string]interface{}{"numbers": []interface{}{1, 2, 3}}})
	t.Logf("%v", res)

	res, err = GetKey("$$.AAA.numbers[5]", map[string]interface{}{"AAA": map[string]interface{}{"numbers": []interface{}{1, 2, 3}}}, nil)
	t.Logf("%v %v", res, err)

	res, err = GetKey("$$$.numbers[0]", nil, map[string]interface{}{"_trigger": map[string]interface{}{"numbers": []interface{}{1, 2, 3}}})
	t.Logf("%v %v", res, err)
}

var dataStr string = `
{
    "store": {
        "book": [
            {
                "category": "reference",
                "author": "Nigel Rees",
                "title": "Sayings of the Century",
                "price": 8.95
            },
            {
                "category": "fiction",
                "author": "Evelyn Waugh",
                "title": "Sword of Honour",
                "price": 12.99
            },
            {
                "category": "fiction",
                "author": "Herman Melville",
                "title": "Moby Dick",
                "isbn": "0-553-21311-3",
                "price": 8.99
            },
            {
                "category": "fiction",
                "author": "J. R. R. Tolkien",
                "title": "The Lord of the Rings",
                "isbn": "0-395-19395-8",
                "price": 22.99
            }
        ],
        "bicycle": {
            "color": "red",
            "price": 19.95
        }
    },
    "expensive": 10
}
`

func parser() {
	//test1()
	//test2()
	//test3()
	test4()
}

func test1() {
	var json_data interface{}
	json.Unmarshal([]byte(dataStr), &json_data)

	res, err := jsonpath.JsonPathLookup(json_data, "$.expensive")
	if err == nil {
		fmt.Println("step 1 res: $.expensive")
		fmt.Println(res)
	}
}

func test2() {
	var json_data interface{}
	json.Unmarshal([]byte(dataStr), &json_data)
	//or reuse lookup pattern
	pat, _ := jsonpath.Compile(`$.store.book[?(@.price < $.expensive)].price`)
	res, _ := pat.Lookup(json_data)
	fmt.Println("step 2 res:")
	fmt.Println(res)
}

func test3() {
	var json_data interface{}
	json.Unmarshal([]byte(dataStr), &json_data)

	res, err := jsonpath.JsonPathLookup(json_data, "$.expensive1")
	if err == nil {
		fmt.Println("step 3 res: $.expensive")
		fmt.Println(res)
	} else {
		fmt.Printf("dddd: %v ", err)
	}
}

func test4() {
	var json_data interface{}
	json.Unmarshal([]byte(dataStr), &json_data)

	res, err := jsonpath.JsonPathLookup(json_data, "$.store.book[0]")
	if err == nil {
		fmt.Println("step 3 res: $.store.book[0]")
		fmt.Println(res)
	} else {
		fmt.Printf("dddd: %v ", err)
	}
}

//jsonpath	result
//$.expensive	10
//$.store.book[0].price	8.95
//$.store.book[-1].isbn	"0-395-19395-8"
//$.store.book[0,1].price	[8.95, 12.99]
//$.store.book[0:2].price	[8.95, 12.99, 8.99]
//$.store.book[?(@.isbn)].price	[8.99, 22.99]
//$.store.book[?(@.price > 10)].title	["Sword of Honour", "The Lord of the Rings"]
//$.store.book[?(@.price <.expensive)].price	[8.95, 8.99]
//$.store.book[:].price	[8.9.5, 12.99, 8.9.9, 22.99]

func TestParser(t *testing.T) {
	parser()
}
