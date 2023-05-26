// Author: huaxr
// Time:   2021/8/2 下午2:32
// Git:    huaxr

package accutil

import (
	_ "unsafe"
)

//go:nosplit
func LockSwap(addr *int32, old, new int32) (swapped bool)
