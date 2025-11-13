package byteops

import (
	"bytes"
	"encoding/binary"
	"math"
)

// Float32ToByteVector converts a slice of float32 to a byte vector
func Float32ToByteVector(input []float32) []byte {
	buf := new(bytes.Buffer)
	for _, val := range input {
		// Convert float32 to IEEE 754 binary representation
		bits := math.Float32bits(val)
		err := binary.Write(buf, binary.LittleEndian, bits)
		if err != nil {
			// Handle error in production code, but for simplicity we'll ignore it here
			continue
		}
	}
	return buf.Bytes()
}
