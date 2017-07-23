package hashmap

import "encoding/binary"

type Value int64

const valueSize = 8

func (v Value) Serialize() []byte {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(v))
	return data
}
func Deserialize(data []byte) Value {
	return Value(binary.LittleEndian.Uint64(data))
}

type TestCase struct {
	key   Key
	value Value
}

type ByKeyValue []TestCase

func (b ByKeyValue) Len() int      { return len(b) }
func (b ByKeyValue) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByKeyValue) Less(i, j int) bool {
	if b[i].key != b[j].key {
		return b[i].value < b[j].value
	}
	return b[i].key < b[j].key
}
