// Author: huaxr
// Time: 2022/7/5 11:02 上午
// Git: huaxr

package ticker

import (
	"context"
	"testing"
	"time"
)

func TestTick(t *testing.T) {
	LaunchJobManager()

	ctx := context.Background()
	tick1 := time.NewTicker(1 * time.Second)
	job1 := func() {
		t.Log("hello")
	}
	a := NewJob(ctx, "tick1", tick1, job1, time.Now().Add(10*time.Second))

	GetManager().Register(a)

	tick2 := time.NewTicker(2 * time.Second)
	job2 := func() {
		t.Log("world")
	}
	b := NewJob(ctx, "tick2", tick2, job2)
	GetManager().Register(b)

	time.Sleep(3 * time.Second)
	GetManager().Revoke("tick2")

	time.Sleep(2 * time.Second)
	tick3 := time.NewTicker(2 * time.Second)
	job3 := func() {
		t.Log("tick3")
	}
	c := NewJob(ctx, "tick3", tick3, job3)
	GetManager().Register(c)
	select {}
}
