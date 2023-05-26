// Author: huaxr
// Time:   2021/10/27 上午11:05
// Git:    huaxr

package flow

type Context interface {
	Get(key string) (val interface{}, exist bool)
	Set(key string, val interface{})
}
