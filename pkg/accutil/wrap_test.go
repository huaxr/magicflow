// Author: huaxr
// Time:   2021/8/4 下午2:29
// Git:    huaxr

package accutil

import (
	"log"
	"testing"
	"time"
)

func f() {
	log.Println("doing something...")
	time.Sleep(1 * time.Second)
	return
}

func f2(a string) {
	log.Println("doing something...", a)
	time.Sleep(1 * time.Second)
	return
}

func TestWrap(t *testing.T) {
	f3 := func() {}
	go Thread("tmp-name3", f3)
	// pay attention to the closure func inside body.
	go Thread("tmp-name", func() { f() })
	go Thread("tmp-name2", func() {
		a := "aaa"
		f2(a)
		return
	})
	time.Sleep(3 * time.Second)
}
