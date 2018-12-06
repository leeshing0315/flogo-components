package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigEndianInt12(t *testing.T) {
	// [6]:195
	// [7]:87
	// [8]:254
	// [9]:143
	// [10]:254
	// [11]:131
	b := []byte{0x20, 0x04}
	result := BigEndianInt12(b)
	assert.Equal(t, int16(-2047), result)
}
