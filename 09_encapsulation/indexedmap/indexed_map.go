package indexedmap

import "sync"

// IndexedMap uses two arrays to store key and value, but
// such that it is possible to retrieve values and key with the index
// for performance issues.
type IndexedMap struct {
	keys   []string
	values []string
}

func NewIndexedMap() *IndexedMap {
	return &IndexedMap{
		keys:   []string{},
		values: []string{},
	}
}

// Add adds the key-value pair into the underlying map.
// Erase the value if it already exists.
func (m *IndexedMap) Add(key, value string) {
	m.keys = append(m.keys, key)
	m.values = append(m.values, value)
}

// Get retrieves the associated values for the given key
func (m *IndexedMap) Get(key string) []string {
	var values []string
	for idx, k := range m.keys {
		if k == key {
			values = append(values, m.values[idx])
		}
	}
	return values
}

type LockedIndexedMap struct {
	mu *sync.Mutex
	*IndexedMap
}

func NewLockedIndexedMap() *LockedIndexedMap {
	return &LockedIndexedMap{
		mu:         &sync.Mutex{},
		IndexedMap: NewIndexedMap(),
	}
}

func (m *LockedIndexedMap) Add(key, value string) {
	m.mu.Lock()
	m.IndexedMap.Add(key, value)
	m.mu.Unlock()
}

func (m *LockedIndexedMap) Get(key string) []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.IndexedMap.Get(key)
}
