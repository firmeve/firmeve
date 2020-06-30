package maps

import "sync"

var (
	mergeLock sync.Mutex
)

func MergeInterface(m1, m2 map[string]interface{}) map[string]interface{} {
	mergeLock.Lock()
	defer mergeLock.Unlock()

	for k, v := range m2 {
		if v2, ok := m2[k]; ok {
			m1[k] = v2
		} else {
			m1[k] = v
		}
	}

	return m1
}
