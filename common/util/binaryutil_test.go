package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigEndianInt12(t *testing.T) {
	b := []byte{0x20, 0x04}
	result := BigEndianInt12(b)
	assert.Equal(t, int16(-2047), result)
}
