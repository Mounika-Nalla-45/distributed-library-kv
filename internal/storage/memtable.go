package storage

import "sort"

type MemTable struct {
	data    map[string]string
	deleted map[string]bool
}

func NewMemTable() *MemTable {
	return &MemTable{
		data:    make(map[string]string),
		deleted: make(map[string]bool),
	}
}

func (m *MemTable) Put(key string, value string) {
	m.data[key] = value
	delete(m.deleted, key)
}

func (m *MemTable) Get(key string) (string, bool) {
	if m.deleted[key] {
		return "", false
	}

	value, ok := m.data[key]
	return value, ok
}

func (m *MemTable) Delete(key string) {
	delete(m.data, key)
	m.deleted[key] = true
}

func (m *MemTable) IsDeleted(key string) bool {
	return m.deleted[key]
}

func (m *MemTable) Keys() []string {
	keys := make([]string, 0, len(m.data))

	for key := range m.data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func (m *MemTable) DeletedKeys() []string {
	keys := make([]string, 0, len(m.deleted))

	for key := range m.deleted {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
