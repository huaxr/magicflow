package core

import (
	"testing"
)

func newTestGlobal(t *testing.T) {
	gCache = new(global)
	gCache.initPlaybooks()
}
