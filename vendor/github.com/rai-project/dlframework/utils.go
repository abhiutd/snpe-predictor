package dlframework

import "golang.org/x/sync/syncmap"

func syncMapKeys(m syncmap.Map) []string {
	keys := []string{}
	m.Range(func(key0 interface{}, value interface{}) bool {
		key, ok := key0.(string)
		if !ok {
			return true
		}
		keys = append(keys, key)
		return true
	})
	return keys
}
