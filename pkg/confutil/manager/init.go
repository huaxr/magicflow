package manager

import (
	"sync"
)

var (
	lock    sync.RWMutex
	dataMap map[string][]byte
)

type LoaderIF interface {
	Read([]string) (map[string][]byte, error)
	Watch(onChange func(data map[string][]byte), configSlice []string)
}

func Init(dir string, configSlice []string, load LoaderIF) {
	if dir == "" {
		panic("config dir not defined yet")
	}
	if load == nil {
		load = NewFileLoader(dir)
	}
	dataMap, err := load.Read(configSlice)
	if err != nil {
		panic(err)
	}
	loadConfig(dataMap)
	go load.Watch(loadConfig, configSlice)
}

func loadConfig(data map[string][]byte) {
	lock.Lock()
	defer lock.Unlock()
	dataMap = data
}

func GetConfigMap() map[string][]byte {
	lock.RLock()
	defer lock.RUnlock()
	return dataMap
}

func GetConfigByKey(key string) []byte {
	lock.RLock()
	defer lock.RUnlock()
	if values, ok := dataMap[key]; ok {
		return values
	}
	return nil
}

func contains(ext string, target []string) bool {
	result := false
	for _, v := range target {
		if ext == v {
			result = true
			break
		}
	}
	return result
}
