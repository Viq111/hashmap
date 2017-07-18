package hashmap

const (
	valueSize = 12
)

// Value is what will be backed into the array.
// You can and should edit its size and method on it to
// get/set your own object
type Value [valueSize]byte
