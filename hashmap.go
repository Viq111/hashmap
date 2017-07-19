/*
Package hashmap is a GC-light implementation when you want to use a map of type:
map[Key][]Value where Key and value DO NOT contains pointers

Usage is:
- Create your map with a sample object you want to use as a value
- Insert an object either with:
  - Insert if you want to insert the normal object
  - UnsafeInsert if you want to avoid the type check
  - InsertSerialized if you already serialized your object in a array of byte
*/
package hashmap

import (
	"encoding/binary"
	"errors"
	"reflect"
	"unsafe"
)

const (
	growthRatio = 2
	loadFactor  = 0.7
)

// Error definitions
var (
	ErrIncorrectType = errors.New("incorrect type")
	ErrZeroAddress   = errors.New("address of zero can't be used")
)

// Hashmap is an array backed hashmap. It will try to be as fast
// for Insert / Lookup / Deletion as possible (hopefully near map implementation)
// But when the Hashmap is GCed, Go only needs to clean the array backing it
type Hashmap struct {
	cellSize  int64
	data      []byte
	cap       int64 // How many keys we can handle
	len       int64 // How many keys we have right now
	valueSize int64
	valueType reflect.Type
}

// NewHashmap returns a hashmap that map a Key to an array of []objSample
// The capacity is actally how many objects you will store, so for example
// if you have 5 keys and all your arrays are of size 10, then you need to set cap to 50
func NewHashmap(objSample interface{}, cap int64) *Hashmap {
	return hashmapFromType(reflect.TypeOf(objSample), int64(unsafe.Sizeof(objSample)), cap)
}

func hashmapFromType(objType reflect.Type, objSize, cap int64) *Hashmap {
	cellSize := int64(keySize + objSize)
	// Since we will grow as soon as we reach cap, we need the slice to be cap * loadFactor
	sliceSize := int(float64(cellSize*cap)*loadFactor) + 1
	data := make([]byte, sliceSize)
	return &Hashmap{
		cellSize:  cellSize,
		data:      data,
		cap:       cap,
		len:       0,
		valueSize: objSize,
		valueType: objType,
	}
}

// Cap returns the capacity of the map before it needs to grow
func (h *Hashmap) Cap() int {
	return int(h.cap)
}

// Len returns the number of values in the map
func (h *Hashmap) Len() int {
	return int(h.len)
}

// Insert an object into the hashmap. If len > cap, it will grow the map
// It also does type checking to be sure your object is the same as what the map
// is supposed to receive
func (h *Hashmap) Insert(key Key, obj interface{}) error {
	// Check that object if the same
	if reflect.TypeOf(obj) != h.valueType {
		return ErrIncorrectType
	}
	return h.UnsafeInsert(key, obj)
}

// UnsafeInsert is exactly the same as Insert but does not do the type checking
func (h *Hashmap) UnsafeInsert(key Key, obj interface{}) error {
	header := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&obj)),
		Len:  int(h.valueSize),
		Cap:  int(h.valueSize),
	}
	data := *(*[]byte)(unsafe.Pointer(&header))
	return h.InsertSerialized(key, data)
}

// InsertSerialized inserts directly the object into the map from a
// slice of bytes. Returns an error if the slice is not of the correct size
func (h *Hashmap) InsertSerialized(key Key, obj []byte) error {
	if key.IsZero() {
		return ErrZeroAddress
	}
	if int64(len(obj)) != h.valueSize {
		return ErrIncorrectType
	}

	if h.len == h.cap {
		h.Grow()
	}
	index, _ := h.getFirstIndex(key)
	index = h.getFirstFree(index)
	sliceSize := int64(len(h.data)) / h.cellSize
	index = index * sliceSize
	binary.LittleEndian.PutUint64(h.data[index:index+keySize], uint64(key.Hash()))
	index += keySize
	copy(h.data[index:], obj)
	return nil
}

// Get returns a slice of objects from the hashmap. It copies the data
// from the map to a separate array so it is safe to use if the Hashmap grows
// (is changed) or destroyed. If dst is provided, it will try to use that slice without
// allocating a new one. If dst is nil or the capacity is not enough, it will create a new one
func (h *Hashmap) Get(dst []interface{}, key Key) []interface{} {
	var objects []interface{}
	index, found := h.getFirstIndex(key)
	if !found {
		return objects // Empty
	}

	keyHash := key.Hash()
	sliceSize := int64(len(h.data)) / h.cellSize
	for { // While we can actually get objects, continue
		cellKey := int64(binary.LittleEndian.Uint64(h.data[index*h.cellSize : index*h.cellSize+keySize]))
		if cellKey == 0 { // No more keys to find
			return objects
		}
		if cellKey == keyHash { // Add object
			value := make([]byte, h.valueSize)
			copy(value, h.data[index*h.cellSize+keySize:index*h.cellSize+keySize+h.valueSize])
			objects = append(objects, value)

		}
		index = (index + 1) % sliceSize
	}
}

// UnsafeGet returns a slice of objects referencing directly to the Hashmap
// it doesn't copy data so as soon as the hashmap grows or is destroyed, you
// will loose the slice!
func (h *Hashmap) UnsafeGet(key Key) []interface{} {
	// ToDo
	return nil
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
	keyHash := key.Hash()
	index := keyHash % sliceSize
	for { // Iterate over the hashmap to find the first key or empty cell
		cellKey := int64(binary.LittleEndian.Uint64(h.data[index*h.cellSize : index*h.cellSize+keySize]))
		if cellKey == keyHash {
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
