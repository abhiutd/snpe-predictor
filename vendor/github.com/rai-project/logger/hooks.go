package logger

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
)

type hooksMap map[string]logrus.Hook
type hooksInfo struct {
	sync.Mutex
	data hooksMap
}

var (
	hooks = hooksInfo{
		data: hooksMap{},
	}
)

func RegisterHook(name string, imp logrus.Hook) {
	hooks.Lock()
	defer hooks.Unlock()
	hooks.data[name] = imp
}

func GetHook(name string) (logrus.Hook, error) {
	hooks.Lock()
	defer hooks.Unlock()
	if h, ok := hooks.data[name]; ok {
		return h, nil
	}
	return nil, errors.Errorf("hook %v not registered", name)
}
