package util

import (
	"errors"
	"fmt"
	"strconv"
	"time"
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

func ParseDateStrFromBCD6(bcd []byte) (string, error) {
	if len(bcd) < 6 {
		return "", errors.New("BCD6 has no enough length")
	}
	value := fmt.Sprintf("20%s-%s-%sT%s:%s:%s+08:00",
		keepTwoDigits(FormatBCD(bcd[0])),
		keepTwoDigits(FormatBCD(bcd[1])),
		keepTwoDigits(FormatBCD(bcd[2])),
		keepTwoDigits(FormatBCD(bcd[3])),
		keepTwoDigits(FormatBCD(bcd[4])),
		keepTwoDigits(FormatBCD(bcd[5])))

	_, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func keepTwoDigits(input string) string {
	if len(input) < 2 {
		return "0" + input
	}
	return input
}
