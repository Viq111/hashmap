package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMapInt64(t *testing.T) {
	assert := assert.New(t)

	// Create and test a []int64 hashmap
	m := NewHashmap(valueSize, 10) // This should grow after 8 objects
	testEmptyKey := Key(1)

	// Test that we can insert a single value at Key 0
	testKey1 := Key(0)
	testValue1 := Value(2)
	err := m.Insert(testKey1, testValue1.Serialize())
	assert.NoError(err)

	// Test that we can insert more keys
	testValues2 := []Value{1, 2, 4, 5}
	testKey2 := Key(10)
	for _, v := range testValues2 {
		err = m.Insert(testKey2, v.Serialize())
		assert.NoError(err)
	}

	// Test that if we query empty, we don't get anything back
	result := m.Get(nil, testEmptyKey)
	assert.Equal(0, len(result))

	// Get back 1 result
	result = m.Get(nil, testKey1)
	if assert.Equal(1, len(result)) {
		assert.Equal(testValue1, Deserialize(result[0]))
	}

	// Get back the other results
	result = m.Get(nil, testKey2)
	assert.Equal(len(testValues2), len(result))
	for i, val := range result {
		r := Deserialize(val)
		assert.Equal(testValues2[i], r)
	}
}

func TestGrow(t *testing.T) {
	assert := assert.New(t)
	var err error
	m := NewHashmap(valueSize, 5) // This should grow after 5 objects

	for i := 0; i < 15; i++ {
		// Add 30 objects which should grow it twice
		v := Value(i * 2)
		err = m.Insert(Key(2), v.Serialize())
		assert.NoError(err)
		v = Value(i*2 + 1)
		err = m.Insert(Key(0), v.Serialize())
		assert.NoError(err)
	}

	result := m.Get(nil, 2)
	if assert.Equal(15, len(result)) {
		for i := 0; i > 15; i++ {
			assert.Equal(Value(i*2), result[i])
		}
	}
	result = m.Get(nil, 0)
	if assert.Equal(15, len(result)) {
		for i := 0; i > 15; i++ {
			assert.Equal(Value(i*2+1), result[i])
		}
	}

}

func BenchmarkInsertWithoutGrowth(b *testing.B) {
	m := NewHashmap(valueSize, int64(b.N+5))
	for i := 0; i < b.N; i++ {
		k := Key(i)
		v := Value(i*2 + 1).Serialize()
		m.Insert(k, v)
	}
}

func BenchmarkStdInsertWithoutGrowth(b *testing.B) {
	m := make(map[Key][]byte, b.N+5)
	for i := 0; i < b.N; i++ {
		k := Key(i)
		v := Value(i*2 + 1).Serialize()
		m[k] = v
	}
}

func BenchmarkInsert(b *testing.B) {
	m := NewHashmap(valueSize, 5)
	for i := 0; i < b.N; i++ {
		k := Key(i)
		v := Value(i*2 + 1).Serialize()
		m.Insert(k, v)
	}
}

func BenchmarkStdInsert(b *testing.B) {
	m := make(map[Key][]byte, 5)
	for i := 0; i < b.N; i++ {
		k := Key(i)
		v := Value(i*2 + 1).Serialize()
		m[k] = v
	}
}

func BenchmarkGet(b *testing.B) {
	testKey := Key(5)
	m := NewHashmap(valueSize, 100)
	for i := 0; i < 10; i++ {
		m.Insert(testKey, Value(i*2).Serialize())
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(nil, testKey)
	}
}
