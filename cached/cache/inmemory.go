package cache

import "sync"

type inMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
	Stat
}

func (m *inMemoryCache) Set(key string, value []byte) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tmpValue, exits := m.c[key]
	if exits {
		m.del(key, tmpValue)
	}
	m.c[key] = value
	m.add(key, value)
	return nil
}

func (m *inMemoryCache) Get(key string) ([]byte, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.c[key], nil
}

func (m *inMemoryCache) Del(key string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tmpValue, exits := m.c[key]
	if exits {
		delete(m.c, key)
		m.del(key, tmpValue)
	}
	return nil
}

func (m *inMemoryCache) GetStat() Stat {
	return m.Stat
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{make(map[string][]byte), sync.RWMutex{}, Stat{}}
}
