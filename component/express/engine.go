// Author: huaxr
// Time:   2021/8/5 上午9:51
// Git:    huaxr

package express

type Engine int

const (
	Js Engine = iota + 1

	Formula

	Valuate
)

type EngineImpl interface {
	EExecute() (bool, error)
}
