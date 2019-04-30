package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"

	"github.com/sigurn/crc16"
)

var serverUri string = "nptcp.cargosmart.com:1024"

// var serverUri string = "itciot-tcp.cargosmart.ai:8080"

// var serverUri string = "ciot-tcp.cargosmart.ai:8080"

// var serverUri string = "172.16.180.204:8033"

var jsonBody string = `
{
	"from": "2019-04-17T16:32:41+08:00",
	"to": "2019-04-30T18:00:00+08:00"
}
`

func main() {
	conn, err := net.Dial("tcp", serverUri)
	if err != nil {
		log.Println(err)
		return
	}

	bufReader := bufio.NewReader(conn)
	bufWriter := bufio.NewWriter(conn)

	replayPacket := genReplayPacket()
	log.Println(replayPacket)
	_, err = bufWriter.Write(replayPacket)
	if err != nil {
		log.Println(err)
		return
	}
	err = bufWriter.Flush()
	if err != nil {
		log.Println(err)
		return
	}

	// ack cmd
	_, err = bufReader.ReadByte()
	if err != nil {
		log.Println(err)
		return
	}
	// ack seqno
	_, err = readCount(bufReader, 2)
	if err != nil {
		log.Println(err)
		return
	}
	// ack dataSegmentLength
	dataSegmentLength, err := readCount(bufReader, 2)
	if err != nil {
		log.Println(err)
		return
	}
	// ack dataSegment
	dataSegment, err := readCount(bufReader, int(binary.BigEndian.Uint16(dataSegmentLength)))
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(dataSegment))
	// ack crc16
	_, err = readCount(bufReader, 2)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func readCount(reader *bufio.Reader, count int) ([]byte, error) {
	var result []byte = make([]byte, count)
	temp := result
	hasRead := 0
	for {
		n, err := reader.Read(temp)
		if err != nil {
			return result, err
		}
		hasRead += n
		if hasRead == count {
			break
		}
		temp = temp[n:]
	}
	return result, nil
}

func genReplayPacket() []byte {
	var buf bytes.Buffer

	// cmd
	buf.WriteByte(0xFC)
	// seqno
	seqNo := make([]byte, 2)
	binary.BigEndian.PutUint16(seqNo, uint16(12345))
	buf.Write(seqNo)
	// dataSegment length & dataSegment
	dataSegment := []byte(jsonBody)
	dataSegmentLen := make([]byte, 2)
	binary.BigEndian.PutUint16(dataSegmentLen, uint16(len(dataSegment)))
	buf.Write(dataSegmentLen)
	buf.Write(dataSegment)
	// crc16
	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(buf.Bytes(), myTable)
	crc := make([]byte, 2)
	binary.BigEndian.PutUint16(crc, checksum)
	buf.Write(crc)

	return buf.Bytes()
}
