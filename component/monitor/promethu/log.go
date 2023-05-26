// Author: XinRui Hua
// Time:   2022/4/15 下午6:11
// Git:    huaxr

package promethu

import "github.com/huaxr/magicflow/component/logx"

type promlog struct {
}

func (p *promlog) Println(v ...interface{}) {
	logx.L().Info(v...)
}
