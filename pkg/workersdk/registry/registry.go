// Author: huaxr
// Time:   2021/8/20 上午11:48
// Git:    huaxr

package registry

import (
	"sync"
)

var (
	registry *Registry
	once     sync.Once
)

type Registry struct {
	sync.RWMutex
	tasks     map[string]Task
	commands  map[string]Command
	app       string
	heartbeat int
}

func init() {
	once.Do(func() {
		registry = &Registry{
			tasks:    make(map[string]Task),
			commands: make(map[string]Command),
		}
	})
}

func RegisterTask(name string, t Task) {
	registry.Lock()
	registry.tasks[name] = t
	registry.Unlock()
}

func GetTask(name string) (Task, bool) {
	registry.RLock()
	defer registry.RUnlock()
	t, ok := registry.tasks[name]
	return t, ok
}

func RevokeTask(name string) {
	registry.Lock()
	delete(registry.tasks, name)
	registry.Unlock()
}

func RegisterCommand(name string, t Command) {
	registry.Lock()
	registry.commands[name] = t
	registry.Unlock()
}

func GetCommand(name string) (Command, bool) {
	registry.RLock()
	defer registry.RUnlock()
	t, ok := registry.commands[name]
	return t, ok
}

func RevokeCommand(name string) {
	registry.Lock()
	delete(registry.commands, name)
	registry.Unlock()
}
