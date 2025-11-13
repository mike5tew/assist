package byteops

import (
	"encoding/binary"
	"math"
)

// Float32ToByteVector converts a float32 to a byte vector
func Float32ToByteVector(f float32) []byte {
	bits := math.Float32bits(f)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
