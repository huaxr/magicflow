// Author: huaxr
// Time:   2021/9/27 下午4:55
// Git:    huaxr

package healthy

import (
	"net/http"
	_ "net/http/pprof"
)

type Profile struct {
}

func (p *Profile) String() string {
	return "Profile"
}

func (p *Profile) Report() {
	_ = http.ListenAndServe("0.0.0.0:10000", nil)
}
