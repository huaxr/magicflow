// Author: XinRui Hua
// Time:   2022/4/24 下午6:02
// Git:    huaxr

package confutil

type Prom struct {
	PullPort    string `yaml:"pullPort"`
	PushGateWay string `yaml:"pushGateWay"`
}

var prom *Prom

func GetProm() *Prom {
	if prom == nil {
		initConf()
	}
	return prom
}

func (l *Prom) GetPullPort() string {
	return l.PullPort
}

func (l *Prom) GetPushGateWay() string {
	return l.PushGateWay
}
