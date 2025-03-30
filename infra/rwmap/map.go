package rwmap

import "sync"

type Map[K comparable, V any] struct {
	m    map[K]V
	lock sync.RWMutex
	once sync.Once
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	m.init()
	m.lock.RLock()
	defer m.lock.RUnlock()
	value, ok := m.m[key]
	return value, ok
}

func (m *Map[K, V]) Set(key K, value V) {
	m.init()
	m.lock.Lock()
	defer m.lock.Unlock()
	m.m[key] = value
}
func (m *Map[K, V]) Delete(key K) {
	m.init()
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.m, key)
}

func (m *Map[K, V]) init() {
	m.once.Do(func() {
		m.lock.Lock()
		defer m.lock.Unlock()
		if m.m == nil {
			m.m = make(map[K]V)
		}
	})
}
