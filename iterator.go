package hashmap

import "encoding/binary"

// Iterator is an iterator over a hashmap
type Iterator struct {
	index     int64
	h         *Hashmap
	zeroIndex int64 // Where we are at for zero values. Goes to -1 when done
}

func newIterator(h *Hashmap) *Iterator {
	return &Iterator{
		h: h,
	}
}

// unsafeNextZero iterate over the values of the zero keys
// returns False if it can't find one anymore
func (i *Iterator) unsafeNextZero(value *[]byte) bool {
	if i.zeroIndex == -1 || i.zeroIndex >= int64(len(i.h.zeroArray)) {
		i.zeroIndex = -1
		return false
	}
	*value = i.h.zeroArray[i.zeroIndex]
	i.zeroIndex++
	return true
}

// Next put the next key (can be the same) and value in the given parameters
// It returns whether the iterator is ended or not
// Usage is:
// for it.Next(&key, &value) {
//   // Do things for key & value
// }
func (i *Iterator) Next(key *Key, value *[]byte) bool {
	var v []byte
	if int64(cap(*value)) >= i.h.valueSize {
		v = (*value)[0:i.h.valueSize]
	} else {
		v = make([]byte, i.h.valueSize)
	}
	var src []byte
	exist := i.UnsafeNext(key, &src)
	if !exist { // End of loop, quit anyway
		return exist
	}
	copy(v, src)
	*value = v
	return true
}

// UnsafeNext is the same as next exact that it does not copy the data into value
// It just returns a pointer to the current in memory structure without copying
func (i *Iterator) UnsafeNext(key *Key, value *[]byte) bool {
	*key = 0
	if i.zeroIndex != -1 { // Still got some zero values
		end := i.unsafeNextZero(value)
		if end {
			return end
		}
	}
	var offset int64
	maxSize := int64(len(i.h.data)) / i.h.cellSize
	for {
		if i.index >= maxSize {
			return false
		}
		offset = i.index * i.h.cellSize
		*key = Key(int64(binary.LittleEndian.Uint64(i.h.data[offset : offset+keySize])))
		if *key == 0 {
			i.index++
		} else {
			break
		}
	}
	offset += keySize
	*value = i.h.data[offset : offset+i.h.valueSize]
	i.index++
	return true

}
