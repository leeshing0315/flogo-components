package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"

	"github.com/sigurn/crc16"
)

// var serverUri string = "nptcp.cargosmart.com:1024"

// var serverUri string = "itciot-tcp.cargosmart.ai:8080"

var serverUri string = "ciot-tcp.cargosmart.ai:8080"

// var serverUri string = "172.16.180.204:8033"

// var jsonBody string = `
// {
//  "from": "2019-04-17T16:32:41+08:00",
//  "to": "2019-04-30T18:00:00+08:00"
// }
// `

var jsonBody string = `
{
    "from": "2019-04-17T16:32:41+08:00",
    "to": "2019-04-19T00:00:00+08:00"
}
`

var jsonBody19 string = `
{
    "from": "2019-04-19T00:00:00+08:00",
    "to": "2019-04-20T00:00:00+08:00"
}
`

var jsonBody20 string = `
{
    "from": "2019-04-20T00:00:00+08:00",
    "to": "2019-04-21T00:00:00+08:00"
}
`

var jsonBody21 string = `
{
    "from": "2019-04-21T00:00:00+08:00",
    "to": "2019-04-22T00:00:00+08:00"
}
`
var jsonBody22 string = `
{
    "from": "2019-04-22T00:00:00+08:00",
    "to": "2019-04-23T00:00:00+08:00"
}
`

var jsonBody23 string = `
{
    "from": "2019-04-23T00:00:00+08:00",
    "to": "2019-04-24T00:00:00+08:00"
}
`

var jsonBody24 string = `
{
    "from": "2019-04-24T00:00:00+08:00",
    "to": "2019-04-25T00:00:00+08:00"
}
`

var jsonBody25 string = `
{
    "from": "2019-04-25T00:00:00+08:00",
    "to": "2019-04-26T00:00:00+08:00"
}
`

var jsonBody26 string = `
{
    "from": "2019-04-26T00:00:00+08:00",
    "to": "2019-04-27T00:00:00+08:00"
}
`

var jsonBody27 string = `
{
    "from": "2019-04-27T00:00:00+08:00",
    "to": "2019-04-28T00:00:00+08:00"
}
`

var jsonBody28 string = `
{
    "from": "2019-04-28T00:00:00+08:00",
    "to": "2019-04-29T00:00:00+08:00"
}
`

var jsonBody29 string = `
{
    "from": "2019-04-29T00:00:00+08:00",
    "to": "2019-04-30T00:00:00+08:00"
}
`

var jsonBodyArr []string = []string{
	jsonBody,
	jsonBody19,
	jsonBody20,
	jsonBody21,
	jsonBody22,
	jsonBody23,
	jsonBody24,
	jsonBody25,
	jsonBody26,
	jsonBody27,
	jsonBody28,
	jsonBody29,
}

var scheduleSeq int = 8

// func main() {
// 	triggerReplay(jsonBody20)
// 	return
// }

func scheduleTask() {
	jsonBodyInput := jsonBodyArr[scheduleSeq]
	go triggerReplay(jsonBodyInput)
	log.Printf("START SCHEDULE TASK %v\n", scheduleSeq)

	scheduleSeq = (scheduleSeq + 1) % len(jsonBodyArr)
	log.Printf("NEXT SCHEDULE TASK WILL BE %v\n", scheduleSeq)
}

func triggerReplay(jsonBodyInput string) {
	conn, err := net.Dial("tcp", serverUri)
	if err != nil {
		log.Println("A " + err.Error())
		return
	}
	defer conn.Close()

	bufReader := bufio.NewReader(conn)
	bufWriter := bufio.NewWriter(conn)

	replayPacket := genReplayPacket(jsonBodyInput)
	log.Println(replayPacket)
	_, err = bufWriter.Write(replayPacket)
	if err != nil {
		log.Println("B " + err.Error())
		return
	}
	err = bufWriter.Flush()
	if err != nil {
		log.Println("C " + err.Error())
		return
	}

	time.Sleep(time.Duration(2) * time.Second)

	// ack cmd
	b, err := bufReader.ReadByte()
	log.Println(b)
	if err != nil {
		log.Println("D " + err.Error())
		return
	}
	// // ack seqno
	// _, err = readCount(bufReader, 2)
	// if err != nil {
	//  log.Println("E " + err.Error())
	//  return
	// }
	// // ack dataSegmentLength
	// dataSegmentLength, err := readCount(bufReader, 2)
	// if err != nil {
	//  log.Println("F " + err.Error())
	//  return
	// }
	// // ack dataSegment
	// dataSegment, err := readCount(bufReader, int(binary.BigEndian.Uint16(dataSegmentLength)))
	// if err != nil {
	//  log.Println("G " + err.Error())
	//  return
	// }
	// log.Println("RESPONSE " + string(dataSegment))
	// // ack crc16
	// _, err = readCount(bufReader, 2)
	// if err != nil {
	//  log.Println("H " + err.Error())
	//  return
	// }
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

func genReplayPacket(jsonBodyInput string) []byte {
	var buf bytes.Buffer

	// pass HTTP filter
	httpLike := []byte{49, 69, 32, 0, 0, 0, 0}
	buf.Write(httpLike)
	// cmd
	buf.WriteByte(0xFC) // 0xFC 252
	// seqno
	// seqNo := make([]byte, 2)
	// binary.BigEndian.PutUint16(seqNo, uint16(12345))
	// buf.Write(seqNo)
	buf.WriteByte(0xff)
	buf.WriteByte(0x32)
	// dataSegment length & dataSegment
	dataSegment := []byte(jsonBodyInput)
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
