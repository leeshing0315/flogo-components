package util

import (
	"encoding/binary"
)

const symbolicBitMask byte = 0x20 // 0010 0000

// BigEndianInt12 parse bit2-bit13
func BigEndianInt12(b []byte) int16 {
	if len(b) < 2 {
		return int16(b[0])
	}
	unsignedNumber := binary.BigEndian.Uint16(b[0:2])
	unsignedNumber <<= 2
	signedNumber := int16(unsignedNumber)
	signedNumber >>= 4
	return signedNumber
}
