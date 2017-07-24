package main

import "encoding/binary"

// Value is a small sample Go object (only value, no pointers)
type Value struct {
	step     int32
	lastSeen int64
}

// Serialize the object into an array of byte
func (v Value) Serialize(buffer []byte) []byte {
	if cap(buffer) >= 12 {
		buffer = buffer[0:12]
	} else {
		buffer = make([]byte, 12)
	}
	binary.LittleEndian.PutUint32(buffer[0:4], uint32(v.step))
	binary.LittleEndian.PutUint64(buffer[4:12], uint64(v.lastSeen))
	return buffer
}

// Deserialize the array to a value
func Deserialize(data []byte) Value {
	return Value{
		step:     int32(binary.LittleEndian.Uint32(data[0:4])),
		lastSeen: int64(binary.LittleEndian.Uint64(data[4:12])),
	}
}
