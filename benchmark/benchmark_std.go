package main

// MapStd is backed by a standard go map
type MapStd struct {
	m map[int64][]Value
}

func NewMapStd(cap int) Map {
	return &MapStd{
		m: make(map[int64][]Value, cap),
	}
}

func (m *MapStd) Insert(key int64, value Value) {
	m.m[key] = append(m.m[key], value)
}

func (m *MapStd) Get(key int64) []Value {
	return m.m[key]
}

func (m *MapStd) Len() int {
	return len(m.m)
}
