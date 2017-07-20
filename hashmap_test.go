package hashmap

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMapInt64(t *testing.T) {
	assert := assert.New(t)

	// Define your value here
	type Value int64
	valueSize := int64(8)
	Serialize := func(v Value) []byte {
		data := make([]byte, 8)
		binary.LittleEndian.PutUint64(data, uint64(v))
		return data
	}
	Deserialize := func(data []byte) Value {
		return Value(binary.LittleEndian.Uint64(data))
	}

	// Create and test a []int64 hashmap
	m := NewHashmap(valueSize, 10) // This should grow after 8 objects
	testEmptyKey := Key(1)

	// Test that we can insert a single value at Key 0
	testKey1 := Key(0)
	testValue1 := Value(2)
	err := m.Insert(testKey1, Serialize(testValue1))
	assert.NoError(err)

	// Test that we can insert more keys
	testValues2 := []Value{1, 2, 4, 5}
	testKey2 := Key(10)
	for _, v := range testValues2 {
		err = m.Insert(testKey2, Serialize(v))
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
