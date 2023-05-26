// Author: huaxr
// Time:   2021/8/3 下午9:35
// Git:    huaxr

package accutil

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContains(t *testing.T) {
	//println(containsNum([]int{1, 2, 3, 4}, 4))
	//println(containsNum([]int{1, 2, 3, 4}, 0))
	println(ContainsStr([]string{"a", "b", "c"}, "a"))
	println(ContainsStr([]string{"a", "b", "c"}, "c"))
	println(ContainsStr([]string{"aaaa", "sssfaagffffffffb", "cweqerfasdadsa", "X", "A", "a", "a", "a", "a", "s", "z"}, "x"))
}

func Test(t *testing.T) {
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("gorutine...")
			runtime.Gosched()
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("main...")
			//runtime.Gosched()
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			fmt.Println("xxx...")
			//runtime.Gosched()
		}
	}()
	time.Sleep(time.Second)
}
