package smusingledata

import (
	"strconv"
)

const (
	firstHalfMask  byte = 0xF0 // 1111 0000
	secondHalfMask byte = 0x0F // 0000 1111
)

func FormatBCD(b byte) string {
	firstHalf := uint8(b) & firstHalfMask >> 4
	secondHalf := b & secondHalfMask
	return strconv.FormatUint(uint64(firstHalf*10+secondHalf), 10)
}
