package main

import (
	"encoding/binary"

	"github.com/sigurn/crc16"
)

func main() {
	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum([]byte{54, 12, 208, 0, 108, 0, 0, 0, 114, 24, 17, 33, 23, 67, 81, 1, 215, 31, 42, 7, 68, 30, 6, 0, 1, 0, 40, 100, 4, 2, 15, 10, 1, 33, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56, 67, 48, 48, 48, 48, 49, 83, 77, 85, 84, 48, 48, 48, 48, 48, 48, 49, 68, 2, 44, 68, 1, 255, 255, 255, 64, 195, 87, 254, 143, 254, 131, 255, 111, 255, 255, 255, 38, 6, 15, 210, 1, 182, 1, 142, 194, 183, 255, 195, 254, 99, 4, 63, 0, 208, 1, 72, 0, 0, 255, 193, 255, 195, 80}, myTable)
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, checksum)
	println(result[0], result[1])
}
