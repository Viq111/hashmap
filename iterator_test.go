package hashmap

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	assert := assert.New(t)
	h := NewHashmap(valueSize, 10)

	testCases := []TestCase{
		TestCase{key: 0, value: 5},
		TestCase{key: 3, value: 7},
		TestCase{key: 2, value: 2},
	}

	for _, test := range testCases {
		err := h.Insert(test.key, test.value.Serialize())
		assert.NoError(err)
	}
	it := newIterator(h)
	var k Key
	var vRaw []byte

	var result []TestCase
	for it.Next(&k, &vRaw) {
		v := Deserialize(vRaw)
		result = append(result, TestCase{key: k, value: v})
	}
	sort.Sort(ByKeyValue(testCases))
	sort.Sort(ByKeyValue(result))
	assert.Equal(testCases, result)
}
