package hashmap

const (
	keySize    = 8
	loadFactor = 0.7
)

// Key denotates what is used as a key in the Hashmap
// In the current implementation, we use an int but you can
// hash your own key to that
type Key int64

// Hashmap is an array backed hashmap. It will try to be as fast
// for Insert / Lookup / Deletion as possible (hopefully near map implementation)
// But when the Hashmap is GCed, Go only needs to clean the array backing it
type Hashmap struct {
	data []byte
	cap  int64 // How many keys we can handle
	len  int64 // How many keys we have right now
}

func NewHashmap(cap int64) *Hashmap {
	cellSize := int64(keySize + valueSize)
	data := make([]byte, cellSize*cap)
	return &Hashmap{
		data: data,
		cap:  cap,
	}
}
