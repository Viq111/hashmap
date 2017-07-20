/*
Package hashmap is a GC-light implementation when you want to use a map of type:
map[Key][]Value where Key and value DO NOT contains pointers
For sake of simplicity and because Go doesn't have generic, Value will be a slice
of bytes. It needs to be fixed/same size for all calls

Usage is:
- Create your map with the vlaue size and capacity
- Insert an object either with:
  - Insert if you want to insert the normal object

// ToDo
*/
package hashmap

import (
	"encoding/binary"
	"errors"
)

const (
	growthRatio = 2
	loadFactor  = 0.7
)

// Error definitions
var (
	ErrIncorrectValueSize = errors.New("incorrect value size")
)

// Hashmap is an array backed hashmap. It will try to be as fast
// for Insert / Lookup / Deletion as possible (hopefully near map implementation)
// But when the Hashmap is GCed, Go only needs to clean the array backing it
type Hashmap struct {
	cellSize  int64
	data      []byte
	capacity  int64 // How many values we can handle
	length    int64 // How many values we have right now
	valueSize int64
	zeroArray [][]byte // Since we can't save the key 0 (used to indicate empty cell), use just array for it
}

// NewHashmap returns a hashmap that map a Key to an array of []objSample
// The capacity is actally how many objects you will store, so for example
// if you have 5 keys and all your arrays are of size 10, then you need to set cap to 50
func NewHashmap(valueSize, capacity int64) *Hashmap {
	cellSize := int64(keySize + valueSize)
	// Since we will grow as soon as we reach cap, we need the slice to be cap * loadFactor
	sliceSize := int(float64(capacity*cellSize)*loadFactor) + 1
	data := make([]byte, sliceSize)
	return &Hashmap{
		cellSize:  cellSize,
		data:      data,
		capacity:  capacity,
		length:    0,
		valueSize: valueSize,
	}
}

// Cap returns the capacity of the map before it needs to grow
func (h *Hashmap) Cap() int {
	return int(h.capacity)
}

// Len returns the number of values in the map
func (h *Hashmap) Len() int {
	return int(h.length)
}

// Insert an object into the hashmap. If len > cap, it will grow the map
// It also does type checking to be sure your object is the same as what the map
// is supposed to receive
func (h *Hashmap) Insert(key Key, value []byte) error {
	// Check that the value is the same size
	if int64(len(value)) != h.valueSize {
		return ErrIncorrectValueSize
	}
	if uint64(key) == 0 {
		v := make([]byte, h.valueSize)
		copy(v, value)
		h.zeroArray = append(h.zeroArray, v)
		return nil
	}

	if h.length == h.capacity {
		h.Grow()
	}
	index, _ := h.getFirstIndex(key)
	index = h.getFirstFree(index)
	sliceIndex := index * h.cellSize
	binary.LittleEndian.PutUint64(h.data[sliceIndex:sliceIndex+keySize], uint64(key))
	sliceIndex += keySize
	copy(h.data[sliceIndex:sliceIndex+h.valueSize], value)
	return nil
}

// Get returns a slice of objects from the hashmap. It copies the data
// from the map to a separate array so it is safe to use if the Hashmap grows
// (is changed) or destroyed. If dst is provided, it will try to use that slice without
// allocating a new one. If dst is nil or the capacity is not enough, it will create a new one
func (h *Hashmap) Get(dst [][]byte, key Key) [][]byte {
	if uint64(key) == 0 {
		if cap(dst) >= len(h.zeroArray) {
			dst = dst[0:len(h.zeroArray)]
		} else {
			dst = make([][]byte, len(h.zeroArray))
		}
		for i, src := range h.zeroArray {
			var d []byte
			if int64(cap(dst[i])) >= h.valueSize {
				d = dst[i][0:h.valueSize]
			} else {
				d = make([]byte, h.valueSize)
			}
			copy(d, src)
			dst[i] = d
		}
		return dst
	}
	index, found := h.getFirstIndex(key)
	if !found {
		return dst[0:0] // Return empty slice, works for nil too
	}

	sliceSize := int64(len(h.data)) / h.cellSize
	dstIndex := int64(0)
	for { // While we can actually get objects, continue
		cellKey := int64(binary.LittleEndian.Uint64(h.data[index*h.cellSize : index*h.cellSize+keySize]))
		if cellKey == 0 { // No more keys to find
			return dst[0:dstIndex]
		}
		if Key(cellKey) == key { // Add object
			var value []byte
			exist := false
			if int64(len(dst)) > dstIndex && int64(len(dst[dstIndex])) == h.valueSize {
				value = dst[dstIndex][0:h.valueSize] // Reuse the buffer
				exist = true
			} else { // Or create a new one
				value = make([]byte, h.valueSize)
			}
			copy(value, h.data[index*h.cellSize+keySize:index*h.cellSize+keySize+h.valueSize])
			if !exist { // Only edit dst slice if we didn't reuse the buffer
				if int64(len(dst)) > dstIndex {
					dst[dstIndex] = value
				} else {
					dst = append(dst, value)
				}
			}
			dstIndex++
		}
		index = (index + 1) % sliceSize
	}
}

// Grow just grows the hashmap once it reaches its capacity
func (h *Hashmap) Grow() {
	//newMap := hashmapFromType(h.valueType, h.valueSize, h.cap*growthRatio)
	// Copy data over
	// ToDo
	//h = newMap
}

// getFirstIndex returns the first index in the array
// where we either found the key or we found an empty space
func (h *Hashmap) getFirstIndex(key Key) (int64, bool) {
	sliceSize := int64(len(h.data)) / h.cellSize
	index := key.Hash() % sliceSize
	for { // Iterate over the hashmap to find the first key or empty cell
		cellKey := int64(binary.LittleEndian.Uint64(h.data[index*h.cellSize : index*h.cellSize+keySize]))
		if Key(cellKey) == key {
			return index, true
		}
		if cellKey == 0 {
			return index, false
		}
		index = (index + 1) % sliceSize
	}
}

// getFirstFree returns either the current index if the key is already 0
// or the first available one
func (h *Hashmap) getFirstFree(index int64) int64 {
	sliceSize := int64(len(h.data)) / h.cellSize
	for { // Iterate over the hashmap to find the first empty cell
		cellKey := int64(binary.LittleEndian.Uint64(h.data[index*h.cellSize : index*h.cellSize+keySize]))
		if cellKey == 0 {
			return index
		}
		index = (index + 1) % sliceSize
	}
}
