package main

import (
	"encoding/binary"
	"net"
)

func main() {
	// conn, err := net.Dial("tcp", "itciot-tcp.cargosmart.ai:8080")
	conn, err := net.Dial("tcp", "localhost:8033")
	// conn, err := net.Dial("tcp", "52.193.135.87:8033")
	if err != nil {
		return
	}
	// sendLoginPacket(conn)
	// sendSinglePacket(conn)
	// sendMultiPacket(conn)
	testSendCmd(conn)
}

func sendLoginPacket(conn net.Conn) {
	// n, err := conn.Write([]byte{50, 12, 207, 0, 27, 2, 15, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56, 0, 8, 72, 83, 49, 56, 49, 49, 48, 56, 17, 112})
	n, err := conn.Write([]byte{80, 32, 84, 32, 73, 32, 73, 32, 80, 32, 80, 32, 13, 10,
		50, 12, 207, 0, 27, 2, 15, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56, 0, 8, 72, 83, 49, 56, 49, 49, 48, 56, 17, 112})

	println(n)
	if err != nil {
		println(err.Error())
		return
	}

	ackBuffer := make([]byte, 7)
	conn.Read(ackBuffer)
	println(ackBuffer)

	shouldBeNothing := make([]byte, 1)
	conn.Read(shouldBeNothing)
	println(shouldBeNothing)
}

func sendSinglePacket(conn net.Conn) {
	// n, err := conn.Write([]byte{54, 57, 166, 0, 108, 0, 0, 0, 114, 24, 17, 33, 23, 67, 81, 1, 215, 31, 42, 7, 68, 30, 6, 0, 1, 0, 40, 100, 4, 2, 15, 10, 1, 33, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56, 67, 48, 48, 48, 48, 49, 83, 77, 85, 84, 48, 48, 48, 48, 48, 48, 49, 68, 2, 44, 68, 1, 255, 255, 255, 64, 195, 87, 254, 143, 254, 131, 255, 111, 255, 255, 255, 38, 6, 15, 210, 1, 182, 1, 142, 194, 183, 255, 195, 254, 99, 4, 63, 0, 208, 1, 72, 0, 0, 255, 193, 255, 195, 80, 246, 254})
	n, err := conn.Write([]byte{80, 32, 84, 32, 73, 32, 73, 32, 80, 32, 80, 32, 13, 10,
		54, 57, 166, 0, 108, 0, 0, 0, 114, 24, 17, 33, 23, 67, 81, 1, 215, 31, 42, 7, 68, 30, 6, 0, 1, 0, 40, 100, 4, 2, 15, 10, 1, 33, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56, 67, 48, 48, 48, 48, 49, 83, 77, 85, 84, 48, 48, 48, 48, 48, 48, 49, 68, 2, 44, 68, 1, 255, 255, 255, 64, 195, 87, 254, 143, 254, 131, 255, 111, 255, 255, 255, 38, 6, 15, 210, 1, 182, 1, 142, 194, 183, 255, 195, 254, 99, 4, 63, 0, 208, 1, 72, 0, 0, 255, 193, 255, 195, 80, 246, 254})

	println(n)
	if err != nil {
		println(err.Error())
		return
	}

	ackBuffer := make([]byte, 7)
	conn.Read(ackBuffer)
	println(ackBuffer)

	shouldBeNothing := make([]byte, 1)
	conn.Read(shouldBeNothing)
	println(shouldBeNothing)
}

func sendMultiPacket(conn net.Conn) {
	n, err := conn.Write([]byte{80, 32, 84, 32, 73, 32, 73, 32, 80, 32, 80, 32, 13, 10,
		55, 57, 177, 3, 104,
		54, 57, 166, 0, 108, 0, 0, 0, 114, 24, 17, 33, 23, 67, 81, 1, 215, 31, 42, 7, 68, 30, 6, 0, 1, 0, 40, 100, 4, 2, 15, 10, 1, 33, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56, 67, 48, 48, 48, 48, 49, 83, 77, 85, 84, 48, 48, 48, 48, 48, 48, 49, 68, 2, 44, 68, 1, 255, 255, 255, 64, 195, 87, 254, 143, 254, 131, 255, 111, 255, 255, 255, 38, 6, 15, 210, 1, 182, 1, 142, 194, 183, 255, 195, 254, 99, 4, 63, 0, 208, 1, 72, 0, 0, 255, 193, 255, 195, 80, 246, 254,
		54, 57, 167, 0, 73, 0, 0, 0, 114, 24, 17, 33, 23, 85, 2, 1, 215, 30, 227, 7, 68, 30, 121, 0, 5, 0, 34, 100, 4, 2, 14, 10, 2, 44, 68, 1, 255, 255, 255, 64, 195, 87, 254, 143, 254, 131, 255, 111, 255, 255, 255, 38, 6, 15, 210, 1, 182, 1, 142, 194, 183, 255, 195, 254, 99, 4, 63, 0, 209, 1, 72, 0, 0, 255, 193, 255, 195, 80, 26, 90,
		54, 57, 168, 0, 73, 0, 0, 0, 114, 24, 17, 33, 24, 16, 1, 1, 215, 31, 95, 7, 68, 30, 28, 0, 9, 0, 162, 100, 4, 2, 15, 10, 2, 44, 68, 1, 255, 255, 255, 64, 195, 87, 254, 143, 254, 123, 255, 111, 255, 255, 255, 38, 6, 15, 204, 1, 182, 1, 142, 194, 183, 255, 195, 254, 99, 4, 63, 0, 209, 1, 72, 0, 0, 255, 192, 255, 195, 80, 175, 78,
		54, 57, 169, 0, 73, 0, 0, 0, 114, 24, 17, 33, 24, 37, 1, 1, 215, 31, 68, 7, 68, 30, 53, 0, 20, 0, 68, 100, 4, 2, 16, 10, 2, 44, 68, 1, 255, 255, 255, 64, 195, 87, 254, 143, 254, 123, 255, 111, 255, 255, 255, 38, 6, 15, 230, 1, 187, 1, 148, 194, 183, 255, 195, 254, 91, 4, 63, 0, 208, 1, 72, 0, 0, 255, 193, 255, 195, 80, 251, 203,
		54, 57, 170, 0, 30, 0, 128, 128, 112, 24, 17, 33, 24, 48, 66, 1, 215, 32, 220, 7, 68, 31, 40, 0, 0, 0, 0, 100, 4, 2, 16, 10, 5, 1, 77, 85, 148,
		54, 57, 171, 0, 73, 0, 0, 0, 114, 24, 17, 33, 24, 64, 1, 1, 215, 31, 184, 7, 68, 30, 204, 0, 20, 0, 31, 100, 4, 2, 16, 10, 2, 44, 68, 1, 255, 255, 255, 77, 195, 87, 254, 143, 254, 123, 255, 111, 255, 255, 255, 38, 6, 15, 223, 1, 187, 1, 142, 194, 183, 255, 195, 254, 91, 4, 63, 0, 208, 1, 72, 0, 0, 255, 193, 255, 195, 80, 71, 153,
		54, 57, 172, 0, 73, 0, 0, 0, 114, 24, 17, 33, 24, 85, 1, 1, 215, 32, 0, 7, 68, 30, 32, 0, 14, 0, 27, 100, 4, 2, 14, 10, 2, 44, 68, 1, 255, 255, 255, 77, 195, 87, 254, 143, 254, 123, 255, 111, 255, 255, 255, 38, 6, 15, 237, 1, 187, 1, 148, 194, 183, 255, 195, 254, 99, 4, 63, 0, 209, 1, 72, 0, 0, 255, 193, 255, 195, 80, 229, 201,
		54, 57, 173, 0, 73, 0, 0, 0, 114, 24, 17, 33, 25, 16, 1, 1, 215, 32, 82, 7, 68, 30, 5, 0, 5, 1, 62, 100, 4, 2, 14, 10, 2, 44, 68, 1, 255, 255, 255, 77, 195, 87, 254, 143, 254, 123, 255, 111, 255, 255, 255, 38, 6, 15, 210, 1, 182, 1, 142, 194, 183, 255, 195, 254, 91, 4, 63, 0, 208, 1, 72, 0, 0, 255, 193, 255, 195, 80, 249, 164,
		54, 57, 174, 0, 73, 0, 0, 0, 114, 24, 17, 33, 25, 37, 1, 1, 215, 33, 151, 7, 68, 29, 251, 0, 33, 0, 0, 100, 4, 2, 14, 10, 2, 44, 68, 1, 255, 255, 255, 77, 195, 87, 254, 143, 254, 131, 255, 111, 255, 255, 255, 38, 6, 15, 223, 1, 182, 1, 142, 194, 183, 255, 195, 254, 99, 4, 63, 0, 209, 1, 72, 0, 0, 255, 192, 255, 195, 80, 114, 163,
		54, 57, 175, 0, 73, 0, 0, 0, 114, 24, 17, 33, 25, 64, 1, 1, 215, 32, 28, 7, 68, 30, 4, 0, 9, 0, 0, 100, 4, 2, 14, 10, 2, 44, 68, 1, 255, 255, 255, 77, 195, 87, 254, 143, 254, 123, 255, 111, 255, 255, 255, 38, 6, 15, 230, 1, 176, 1, 137, 194, 183, 255, 195, 254, 91, 4, 63, 0, 209, 1, 72, 0, 0, 255, 192, 255, 195, 80, 189, 218,
		54, 57, 176, 0, 73, 0, 0, 0, 114, 24, 17, 33, 25, 85, 1, 1, 215, 32, 3, 7, 68, 29, 74, 0, 24, 0, 7, 100, 4, 2, 15, 10, 2, 44, 68, 1, 255, 255, 255, 77, 195, 87, 254, 143, 254, 123, 255, 111, 255, 255, 255, 38, 6, 15, 243, 1, 182, 1, 137, 194, 183, 255, 195, 254, 99, 4, 63, 0, 209, 1, 72, 0, 0, 255, 193, 255, 195, 80, 231, 214,
		162, 56,
	})

	println(n)
	if err != nil {
		println(err.Error())
		return
	}

	ackBuffer := make([]byte, 7)
	conn.Read(ackBuffer)
	println(ackBuffer)

	shouldBeNothing := make([]byte, 1)
	conn.Read(shouldBeNothing)
	println(shouldBeNothing)
}

func testSendCmd(conn net.Conn) {
	// send login packet (first send devid)
	n, err := conn.Write([]byte{80, 32, 84, 32, 73, 32, 73, 32, 80, 32, 80, 32, 13, 10,
		50, 12, 207, 0, 27, 2, 15, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56,
		6, 67, 48, 51, 54, 50, 53, // C03625
		8, 72, 83, 49, 56, 49, 49, 48, 56, 17, 112})
	println(n)
	if err != nil {
		println(err.Error())
		return
	}
	loginAckBuffer := make([]byte, 7)
	conn.Read(loginAckBuffer)

	// send data packet
	n, err = conn.Write([]byte{80, 32, 84, 32, 73, 32, 73, 32, 80, 32, 80, 32, 13, 10,
		54, 12, 208, 0, 108, 0, 0, 0, 114, 24, 17, 33, 23, 67, 81, 1, 215, 31, 42, 7, 68, 30, 6, 0, 1, 0, 40, 100, 4, 2, 15, 10, 1, 33, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56, 67, 48, 48, 48, 48, 49, 83, 77, 85, 84, 48, 48, 48, 48, 48, 48, 49, 68, 2, 44, 68, 1, 255, 255, 255, 64, 195, 87, 254, 143, 254, 131, 255, 111, 255, 255, 255, 38, 6, 15, 210, 1, 182, 1, 142, 194, 183, 255, 195, 254, 99, 4, 63, 0, 208, 1, 72, 0, 0, 255, 193, 255, 195, 80, 246, 254})
	println(n)
	if err != nil {
		println(err.Error())
		return
	}
	dataAckBuffer := make([]byte, 7)
	conn.Read(dataAckBuffer)

	// receive cmd
	setBuffer, err := parseByProtocol(conn)
	if err != nil {
		println(err.Error())
		return
	}
	println(setBuffer)

	// send cmd ack packet
	n, err = conn.Write([]byte{80, 32, 84, 32, 73, 32, 73, 32, 80, 32, 80, 32, 13, 10,
		54, 12, 209, 0, 0, 1, 1})
	if err != nil {
		println(err.Error())
		return
	}

	ackBuffer := make([]byte, 7)
	conn.Read(ackBuffer)
	println(ackBuffer)
}

func parseByProtocol(conn net.Conn) ([]byte, error) {
	command := make([]byte, 1)
	_, err := conn.Read(command)
	if err != nil {
		return nil, err
	}

	sequence := make([]byte, 2)
	_, err = conn.Read(sequence)
	if err != nil {
		return nil, err
	}

	dataSegmentLength := make([]byte, 2)
	_, err = conn.Read(dataSegmentLength)
	if err != nil {
		return nil, err
	}

	dataSegment := make([]byte, binary.BigEndian.Uint16(dataSegmentLength))
	_, err = conn.Read(dataSegment)
	if err != nil {
		return nil, err
	}

	crc16Check := make([]byte, 2)
	_, err = conn.Read(crc16Check)
	if err != nil {
		return nil, err
	}

	result := []byte{}
	result = append(result, command...)
	result = append(result, sequence...)
	result = append(result, dataSegmentLength...)
	result = append(result, dataSegment...)
	result = append(result, crc16Check...)
	return result, nil
}
