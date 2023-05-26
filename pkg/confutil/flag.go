// Author: XinRui Hua
// Time:   2022/3/17 下午4:03
// Git:    huaxr

package confutil

import (
	"flag"
	"runtime"
)

var (
	confDir *string
)

func init() {
	sysType := runtime.GOOS
	if sysType == "darwin" {
		confDir = flag.String("dir", "/Users/huaxinrui/go/src/github.com/huaxr/magicflow/conf/local", "config yml dirs")
		return
	}

	confDir = flag.String("dir", "", "config yml dirs")
	if !flag.Parsed() {
		flag.Parse()
	}

	if len(*confDir) == 0 {
		panic("flag dir not init yet")
	}
}
