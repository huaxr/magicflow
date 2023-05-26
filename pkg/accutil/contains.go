// Author: huaxr
// Time:   2021/8/3 下午9:34
// Git:    huaxr

package accutil

//go:nosplit
func containsNum([]int, int) int

//go:nosplit
func containsStr(a []string, b string) int

// ContainsNum using to check if b exist in a
func ContainsNumX(a []int, b int) bool {
	if containsNum(a, b) == 1 {
		return true
	}
	return false
}

// bugs need repair.
func ContainsStrX(a []string, b string) bool {
	res := containsStr(a, b)
	if res == 1 {
		return true
	}
	return false
}

func ContainsStr(s []string, val string) bool {
	for _, i := range s {
		if i == val {
			return true
		}
	}
	return false
}

func ContainsNum(a []int, b int) bool {
	for _, i := range a {
		if i == b {
			return true
		}
	}
	return false
}
