package config

import "sync"

var (
	registryMutex sync.Mutex
	registry      = []ConfigInterface{}
)

// Register ...
func Register(mod ConfigInterface) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	registry = append(registry, mod)
}
