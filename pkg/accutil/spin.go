// Author: huaxr
// Time:   2021/8/7 下午11:23
// Git:    huaxr

package accutil

// The spin lock will not switch the thread state. It is always in the user state, that is,
// the thread is permanent active; The thread will not enter the blocking state,
// unnecessary context switching is reduced, and the execution speed is fast
func procyield(cycles uint32)
