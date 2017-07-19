package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMapInt64(t *testing.T) {
	// Create and test a []int64 hashmap
	assert := assert.New(t)
	i := int64(5)
	m := NewHashmap(i, 50)

	sample := []int64{1, 2, 4, 5}
	testKey := Key(10)
	for _, s := range sample {
		err := m.Insert(testKey, s)
		assert.NoError(err)
	}

	result := m.Get(nil, testKey)
	t.Logf("%v", m.data)
	assert.Equal(len(sample), len(result))
	for j, val := range result {
		v, ok := val.(*int64)
		assert.True(ok)
		assert.Equal(sample[j], *v)
	}
}
