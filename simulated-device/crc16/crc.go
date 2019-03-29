package main

import (
	"bytes"
	"encoding/binary"

	"github.com/sigurn/crc16"
)

var pin string = "460010604706821" // C01937 (15 bytes)

var firmwareVersion string = "HS19-2"

func main() {
	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	// checksum := crc16.Checksum([]byte{51, 0, 0, 0, 124, 42, 81, 0, 0, 118, 101, 101, 106, 116, 106, 52, 53, 57, 0, 0, 224, 139, 251, 172, 72, 83, 49, 57, 45, 51, 32, 32, 49, 57, 45, 48, 51, 45, 50, 53, 48, 57, 58, 53, 52, 58, 53, 56, 72, 83, 49, 57, 45, 51, 32, 32, 105, 110, 116, 101, 114, 110, 97, 108, 38, 108, 116, 59, 112, 38, 103, 116, 59, 84, 104, 105, 115, 32, 105, 115, 32, 72, 83, 49, 57, 45, 51, 46, 38, 108, 116, 59, 47, 112, 38, 103, 116, 59, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 35}, myTable)
	checksum := crc16.Checksum([]byte{52, 0, 0, 0, 0}, myTable)
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, checksum)
	println(result[0], result[1])

	// // 示例1。
	// a1 := [7]int{1, 2, 3, 4, 5, 6, 7}
	// fmt.Printf("a1: %v (len: %d, cap: %d)\n",
	// 	a1, len(a1), cap(a1))
	// s9 := a1[1:4]
	// //s9[0] = 1
	// fmt.Printf("s9: %v (len: %d, cap: %d)\n",
	// 	s9, len(s9), cap(s9))
	// for i := 1; i <= 5; i++ {
	// 	s9 = append(s9, i)
	// 	fmt.Printf("s9(%d): %v (len: %d, cap: %d)\n",
	// 		i, s9, len(s9), cap(s9))
	// }
	// fmt.Printf("a1: %v (len: %d, cap: %d)\n",
	// 	a1, len(a1), cap(a1))
	// fmt.Println()

	// loginPacket := genLoginPacket()
	// testPacket := genPacket(0x31, []byte{0, 232}, []byte{})
	// var result []byte
	// result = append(result, loginPacket...)
	// result = append(result, testPacket...)
	// result = append(result, testPacket...)
	// result = append(result, testPacket...)
	// result = append(result, testPacket...)
	// encoded := base64.StdEncoding.EncodeToString(result)
	// log.Println(encoded)
}

func genPacket(cmd byte, seqnoBytes []byte, dataSegment []byte) []byte {
	var packet bytes.Buffer
	packet.WriteByte(cmd)
	packet.Write(seqnoBytes)
	dataSegmentLength := make([]byte, 2)
	binary.BigEndian.PutUint16(dataSegmentLength, uint16(len(dataSegment)))
	packet.Write(dataSegmentLength)
	packet.Write(dataSegment)
	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(dataSegment, myTable)
	crc := make([]byte, 2)
	binary.LittleEndian.PutUint16(crc, checksum)
	packet.Write(crc)
	return packet.Bytes()
}

func genLoginPacket() []byte {
	var loginDataSegmentBuf bytes.Buffer
	loginDataSegmentBuf.WriteByte(0x02)
	loginDataSegmentBuf.WriteByte(byte(len(pin)))
	loginDataSegmentBuf.WriteString(pin)
	loginDataSegmentBuf.WriteByte(0)
	loginDataSegmentBuf.WriteByte(byte(len(firmwareVersion)))
	loginDataSegmentBuf.WriteString(firmwareVersion)
	loginDataSegment := loginDataSegmentBuf.Bytes()
	var loginPacket bytes.Buffer
	loginPacket.WriteByte(0x32)
	loginPacket.Write([]byte{1, 219})
	loginPacket.WriteByte(0)
	loginPacket.WriteByte(byte(len(loginDataSegment)))
	loginPacket.Write(loginDataSegment)
	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(loginPacket.Bytes(), myTable)
	crc := make([]byte, 2)
	binary.LittleEndian.PutUint16(crc, checksum)
	loginPacket.Write(crc)
	return loginPacket.Bytes()
}
