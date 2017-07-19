package hashmap

import "unsafe"

const (
	keySize = 8
)

// Key denotates what is used as a key in the Hashmap
// In the current implementation, we use an int64 but you can
// hash your own key to that. (Make sure if you edit the key type to also edit the keySize)
type Key int64

// Hash should provide a fast hashing of the key to an int64
func (k Key) Hash() int64 {
	// splitmix64 (http://xorshift.di.unimi.it/splitmix64.c)
	// is a good hash function
	i := int64(k)
	x := *(*uint64)(unsafe.Pointer(&i))
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)
	return *(*int64)(unsafe.Pointer(&x))
}

// IsZero returns whether the key is zero
func (k Key) IsZero() bool {
	return k == 0
}
