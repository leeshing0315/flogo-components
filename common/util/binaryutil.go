package util

import (
	"encoding/binary"
	"strconv"
	"strings"
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

func BigEndianFromBit0ToBit6(b byte) int64 {
	b <<= 1
	b >>= 2
	return int64(b)
}

func FromStrToUint32(str string) []byte {
	value, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return nil
	}
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(value))
	return b
}

func FromDecStrToHexStr(str string) []byte {
	value, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return nil
	}
	return []byte(strings.ToUpper(strconv.FormatUint(value, 16)))
}

func GetEndBytes(input []byte, size int) []byte {
	inputLen := len(input)
	result := make([]byte, size)
	if inputLen < size {
		copy(result[size-inputLen:], input)
	} else {
		copy(result, input[inputLen-size:])
	}
	return result
}
