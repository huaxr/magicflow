// Author: huaxr
// Time: 2022/7/4 11:27 上午
// Git: huaxr

package accutil

import (
	_ "unsafe" // for go:linkname
)

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

// AbsTime represents absolute monotonic time.
type AbsTime int64

// Now returns the current absolute monotonic time.
func Now() AbsTime {
	return AbsTime(nanotime())
}
