package main

import "github.com/Viq111/hashmap"

// MapHashmap is backed by a standard go map
type MapHashmap struct {
	m *hashmap.Hashmap
}

func NewMapHashmap(cap int) Map {
	return &MapHashmap{
		m: hashmap.NewHashmap(8, int64(cap)),
	}
}

func (m *MapHashmap) Insert(key int64, value Value) {
	v := value.Serialize(nil)
	m.m.Insert(hashmap.Key(key), v)
}

func (m *MapHashmap) Get(key int64) []Value {
	values := m.m.Get(nil, hashmap.Key(key))

	ret := make([]Value, len(values))
	for i := range values {
		ret[i] = Deserialize(values[i])
	}

	return ret
}

func (m *MapHashmap) Len() int {
	return m.m.Len()
}
