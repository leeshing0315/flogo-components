package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"os"

	"github.com/sigurn/crc16"
)

func sendPacket(writer *bufio.Writer, errChain chan error) {
	stdinReader := bufio.NewReader(os.Stdin)
	for lineBytes, _, err := stdinReader.ReadLine(); err == nil; {
		line := string(lineBytes)
		if line == "q" {
			err = toSendSinglePacket(writer)
			if err != nil {
				errChain <- err
				return
			}
		} else if line == "w" {
			err = toSendMultiPacket(writer)
			if err != nil {
				errChain <- err
				return
			}
		} else if line == "e" {
			err = toSendEventLog(writer)
			if err != nil {
				errChain <- err
				return
			}
		} else if line == "1" {
			err = toSendTest(writer)
			if err != nil {
				errChain <- err
				return
			}
		} else if line == "2" {
			err = toSendTest2Times(writer)
			if err != nil {
				errChain <- err
			}
		} else if line == "3" {
			err = toSendTest3Times(writer)
			if err != nil {
				errChain <- err
			}
		} else if line == "4" {
			err = toSendTest4Times(writer)
			if err != nil {
				errChain <- err
			}
		} else if line == "5" {
			err = toSendTest5Times(writer)
			if err != nil {
				errChain <- err
			}
		} else if line == "p" {
			err = toSendSomethingElse(writer)
			if err != nil {
				errChain <- err
			}
		} else if line == "t" {
			err = toSendSinglePacketWithAutoReg(writer)
			if err != nil {
				errChain <- err
				return
			}
		} else if line == "y" {
			err = toSendMultiPacketWithAutoReg(writer)
			if err != nil {
				errChain <- err
				return
			}
		} else {
			log.Println("******")
		}
		line, err = stdinReader.ReadString('\n')
	}
	errChain <- errors.New("Control-C")
}

func toSendSinglePacket(writer *bufio.Writer) error {
	log.Println("toSendSinglePacket Start")
	dataSegment := []byte{
		0, 128, 0, 114, 25, 3, 5, 16, 7, 6, 1, 222, 173, 70, 7, 63, 30, 104, 0, 0, 0, 247, 100,
		4, 2, 17, 15,
		2, 44, 68, 1, 255, 255, 255, 195, 198, 67, 198, 87, 197, 239, 160, 54, 255, 255, 255, 38, 29, 14, 236, 3, 217, 2, 72, 194, 239, 2, 50, 198, 119, 3, 57, 0, 196, 0, 40, 0, 0, 0, 252, 0, 238, 65,
		8, 2, 81, 93,
		9, 8, 254, 47, 254, 47, 254, 47, 254, 47,
		5, 1, 195,
	}
	packet := genPacket(0x36, []byte{0, 230}, dataSegment)
	_, err := writer.Write(packet)
	if err != nil {
		return err
	}
	log.Println("toSendSinglePacket Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendSinglePacket Flush")
	return nil
}

func toSendMultiPacket(writer *bufio.Writer) error {
	log.Println("toSendMultiPacket Start")
	dataSegment := []byte{
		54,
		27, 79,
		0, 90,
		0, 128, 0, 114, 25, 3, 5, 9, 72, 85, 1, 71, 213, 43, 2, 85, 110, 114, 0, 0, 0, 110, 100,
		4, 2, 18, 17,
		2, 44, 68, 1, 255, 255, 255, 67, 251, 3, 249, 175, 251, 3, 99, 68, 255, 255, 255, 37, 8, 17, 173, 3, 178, 2, 189, 197, 87, 254, 215, 249, 35, 2, 205, 255, 150, 1, 72, 0, 58, 254, 255, 255, 56, 65,
		8, 2, 192, 109,
		9, 8, 254, 47, 254, 47, 254, 47, 254, 47,
		5, 1, 67,
		40, 213,
		54,
		27, 80,
		0, 90,
		0, 128, 0, 114, 25, 3, 5, 9, 87, 66, 1, 71, 212, 245, 2, 85, 110, 52, 0, 0, 0, 110, 100,
		4, 2, 18, 16,
		2, 44, 68, 1, 255, 255, 255, 65, 251, 3, 251, 67, 251, 75, 137, 82, 255, 255, 255, 37, 8, 17, 192, 0, 134, 0, 0, 197, 119, 255, 58, 251, 47, 1, 130, 0, 5, 1, 72, 0, 0, 255, 66, 255, 66, 65,
		8, 2, 0, 36,
		9, 8, 254, 47, 254, 47, 254, 47, 254, 47,
		5, 1, 65,
		39, 175,
		54,
		27, 81,
		0, 90,
		0, 128, 0, 114, 25, 3, 5, 16, 5, 3, 1, 71, 212, 247, 2, 85, 109, 243, 0, 0, 0, 110, 100,
		4, 2, 18, 16,
		2, 44, 68, 1, 255, 255, 255, 67, 251, 3, 249, 183, 251, 3, 102, 68, 255, 255, 255, 37, 8, 17, 173, 3, 178, 2, 184, 197, 75, 254, 218, 249, 11, 2, 222, 255, 90, 1, 72, 0, 58, 255, 2, 255, 56, 65,
		8, 2, 192, 109,
		9, 8, 254, 47, 254, 47, 254, 47, 254, 47,
		5, 1, 67,
		41, 134,
	}
	packet := genPacket(0x37, []byte{27, 82}, dataSegment)
	_, err := writer.Write(packet)
	if err != nil {
		return err
	}
	log.Println("toSendMultiPacket Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendMultiPacket Flush")
	return nil
}

func toSendEventLog(writer *bufio.Writer) error {
	log.Println("toSendEventLog Start")
	dataSegment := []byte{
		1, 19, 145, 163, 128, 56, 238, 52, 71, 72, 222, 52, 47, 60, 193, 255, 255, 255, 255, 255, 1, 19, 145, 163, 192, 56, 222, 51, 71, 72, 222, 52, 47, 61, 193, 255, 255, 255, 255, 255, 1, 19, 145, 164, 0, 56, 234, 52, 71, 72, 222, 52, 47, 61, 193, 255, 255, 255, 255, 255, 1, 19, 145, 164, 64, 56, 230, 52, 71, 72, 226, 52, 47, 61, 193, 255, 255, 255, 255, 255, 1, 19, 145, 164, 128, 56, 230, 52, 71, 72, 226, 52, 48, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 164, 192, 56, 226, 51, 71, 72, 230, 54, 48, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 165, 0, 56, 230, 51, 71, 72, 222, 51, 48, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 165, 64, 56, 226, 51, 71, 72, 234, 54, 48, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 165, 128, 56, 234, 52, 71, 72, 218, 51, 50, 188, 193, 255, 255, 255, 255, 255, 1, 19, 145, 168, 64, 232, 114, 35, 71, 72, 110, 36, 47, 62, 193, 255, 255, 255, 255, 255, 1, 19, 145, 165, 192, 56, 222, 50, 71, 72, 230, 53, 49, 189, 193, 255, 255, 255, 255, 255, 2, 19, 145, 168, 89, 160, 2, 2, 19, 145, 168, 89, 160, 1, 2, 19, 145, 165, 192, 32, 1, 1, 19, 145, 168, 128, 217, 94, 65, 71, 72, 222, 45, 51, 57, 193, 255, 255, 255, 255, 255, 1, 19, 145, 168, 192, 56, 226, 52, 71, 73, 6, 61, 50, 61, 193, 255, 255, 255, 255, 255, 1, 19, 145, 169, 0, 56, 230, 52, 71, 72, 222, 52, 49, 61, 193, 255, 255, 255, 255, 255, 1, 19, 145, 169, 64, 56, 238, 53, 71, 72, 222, 52, 49, 60, 193, 255, 255, 255, 255, 255, 1, 19, 145, 169, 128, 56, 222, 51, 71, 72, 222, 52, 49, 61, 193, 255, 255, 255, 255, 255, 1, 19, 145, 169, 192, 56, 234, 52, 71, 72, 222, 52, 48, 60, 193, 255, 255, 255, 255, 255, 1, 19, 145, 170, 0, 56, 226, 51, 71, 72, 222, 52, 48, 61, 193, 255, 255, 255, 255, 255, 1, 19, 145, 170, 64, 56, 226, 51, 71, 72, 222, 52, 48, 60, 193, 255, 255, 255, 255, 255, 1, 19, 145, 170, 128, 56, 238, 53, 71, 72, 222, 52, 48, 60, 193, 255, 255, 255, 255, 255, 1, 19, 145, 170, 192, 56, 222, 51, 71, 72, 222, 52, 47, 61, 193, 255, 255, 255, 255, 255, 1, 19, 145, 171, 0, 56, 238, 52, 71, 72, 222, 52, 47, 60, 193, 255, 255, 255, 255, 255, 1, 19, 145, 171, 64, 56, 226, 52, 71, 72, 226, 52, 47, 61, 193, 255, 255, 255, 255, 255, 1, 19, 145, 171, 128, 56, 222, 51, 71, 72, 230, 53, 48, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 171, 192, 56, 230, 52, 71, 72, 230, 53, 48, 188, 193, 255, 255, 255, 255, 255, 1, 19, 145, 172, 0, 56, 222, 50, 71, 72, 214, 51, 47, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 172, 64, 56, 226, 51, 71, 72, 226, 52, 47, 188, 193, 255, 255, 255, 255, 255, 1, 19, 145, 172, 128, 56, 222, 50, 71, 72, 214, 51, 47, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 172, 192, 56, 226, 51, 71, 72, 226, 52, 47, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 173, 0, 56, 226, 51, 71, 72, 218, 50, 47, 188, 193, 255, 255, 255, 255, 255, 1, 19, 145, 173, 64, 56, 222, 50, 71, 72, 230, 52, 47, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 173, 128, 56, 222, 50, 71, 72, 218, 51, 46, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 173, 192, 56, 222, 50, 71, 72, 226, 52, 46, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 176, 0, 56, 226, 51, 71, 72, 222, 51, 47, 188, 193, 255, 255, 255, 255, 255, 1, 19, 145, 176, 64, 56, 226, 51, 71, 72, 226, 52, 48, 189, 193, 255, 255, 255, 255, 255, 1, 19, 145, 176, 128, 56, 230, 52, 71, 72, 218, 51, 49, 188, 193, 255, 255, 255, 255, 255, 1, 19, 145, 176, 192, 40, 226, 51, 71, 72, 226, 52, 49, 188, 193, 255, 255, 255, 255, 255, 1, 19, 145, 177, 0, 56, 234, 52, 71, 72, 226, 53, 50, 188, 193, 255, 255, 255, 255, 255, 1, 19, 145, 177, 64, 56, 222, 51, 71, 72, 218, 52, 50, 188, 193, 255, 255, 255, 255, 255,
	}
	packet := genPacket(0x38, []byte{0, 231}, dataSegment)
	_, err := writer.Write(packet)
	if err != nil {
		return err
	}
	log.Println("toSendEventLog Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendEventLog Flush")
	return nil
}

func toSendTest(writer *bufio.Writer) error {
	log.Println("toSendTest Start")
	packet := genPacket(0x31, []byte{0, 232}, []byte{})
	_, err := writer.Write(packet)
	if err != nil {
		return err
	}
	log.Println("toSendTest Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendTest Flush")
	return nil
}

func toSendTest2Times(writer *bufio.Writer) error {
	log.Println("toSendTest2Times Start")
	packet := genPacket(0x31, []byte{0, 232}, []byte{})
	var bigPacket []byte
	bigPacket = append(bigPacket, packet...)
	bigPacket = append(bigPacket, packet...)
	log.Println(bigPacket)
	_, err := writer.Write(bigPacket)
	if err != nil {
		return err
	}
	log.Println("toSendTest2Times Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendTest2Times Flush")
	return nil
}

func toSendTest3Times(writer *bufio.Writer) error {
	log.Println("toSendTest3Times Start")
	packet := genPacket(0x31, []byte{0, 232}, []byte{})
	var bigPacket []byte
	bigPacket = append(bigPacket, packet...)
	bigPacket = append(bigPacket, packet...)
	bigPacket = append(bigPacket, packet...)
	log.Println(bigPacket)
	_, err := writer.Write(bigPacket)
	if err != nil {
		return err
	}
	log.Println("toSendTest3Times Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendTest3Times Flush")
	return nil
}

func toSendTest4Times(writer *bufio.Writer) error {
	log.Println("toSendTest4Times Start")
	packet := genPacket(0x31, []byte{0, 232}, []byte{})
	var bigPacket []byte
	bigPacket = append(bigPacket, packet...)
	bigPacket = append(bigPacket, packet...)
	bigPacket = append(bigPacket, packet...)
	bigPacket = append(bigPacket, packet...)
	log.Println(bigPacket)
	_, err := writer.Write(bigPacket)
	if err != nil {
		return err
	}
	log.Println("toSendTest4Times Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendTest4Times Flush")
	return nil
}

func toSendTest5Times(writer *bufio.Writer) error {
	log.Println("toSendTest5Times Start")
	packet := genPacket(0x31, []byte{0, 232}, []byte{})
	var bigPacket []byte
	bigPacket = append(bigPacket, packet...)
	bigPacket = append(bigPacket, packet...)
	bigPacket = append(bigPacket, packet...)
	bigPacket = append(bigPacket, packet...)
	bigPacket = append(bigPacket, packet...)
	log.Println(bigPacket)
	_, err := writer.Write(bigPacket)
	if err != nil {
		return err
	}
	log.Println("toSendTest5Times Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendTest5Times Flush")
	return nil
}

func toSendSomethingElse(writer *bufio.Writer) error {
	log.Println("toSendSomethingElse Start")
	packet := genPacket(0x99, []byte{0, 232}, []byte{})
	log.Println(packet)
	_, err := writer.Write(packet)
	if err != nil {
		return err
	}
	log.Println("toSendSomethingElse Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendSomethingElse Flush")
	return nil
}

func toSendSinglePacketWithAutoReg(writer *bufio.Writer) error {
	log.Println("toSendSinglePacketWithAutoReg Start")
	dataSegment := []byte{
		0, 128, 0, 114, 25, 3, 5, 16, 7, 6, 1, 222, 173, 70, 7, 63, 30, 104, 0, 0, 0, 247, 100,
		4, 2, 17, 15,
		2, 44, 68, 1, 255, 255, 255, 195, 198, 67, 198, 87, 197, 239, 160, 54, 255, 255, 255, 38, 29, 14, 236, 3, 217, 2, 72, 194, 239, 2, 50, 198, 119, 3, 57, 0, 196, 0, 40, 0, 0, 0, 252, 0, 238, 65,
		8, 2, 81, 93,
		9, 8, 254, 47, 254, 47, 254, 47, 254, 47,
		5, 1, 195,
	}
	dataSegment = appendLoginItem(dataSegment)
	packet := genPacket(0x36, []byte{0, 230}, dataSegment)
	_, err := writer.Write(packet)
	if err != nil {
		return err
	}
	log.Println("toSendSinglePacketWithAutoReg Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendSinglePacketWithAutoReg Flush")
	return nil
}

func toSendMultiPacketWithAutoReg(writer *bufio.Writer) error {
	log.Println("toSendMultiPacketWithAutoReg Start")
	firstSubDataSegment := []byte{
		0, 128, 0, 114, 25, 3, 5, 9, 72, 85, 1, 71, 213, 43, 2, 85, 110, 114, 0, 0, 0, 110, 100,
		4, 2, 18, 17,
		2, 44, 68, 1, 255, 255, 255, 67, 251, 3, 249, 175, 251, 3, 99, 68, 255, 255, 255, 37, 8, 17, 173, 3, 178, 2, 189, 197, 87, 254, 215, 249, 35, 2, 205, 255, 150, 1, 72, 0, 58, 254, 255, 255, 56, 65,
		8, 2, 192, 109,
		9, 8, 254, 47, 254, 47, 254, 47, 254, 47,
		5, 1, 67,
	}
	firstSubDataSegment = appendLoginItem(firstSubDataSegment)
	firstSubPacket := genPacket(54, []byte{27, 79}, firstSubDataSegment)

	secondSubDataSegment := []byte{
		0, 128, 0, 114, 25, 3, 5, 9, 87, 66, 1, 71, 212, 245, 2, 85, 110, 52, 0, 0, 0, 110, 100,
		4, 2, 18, 16,
		2, 44, 68, 1, 255, 255, 255, 65, 251, 3, 251, 67, 251, 75, 137, 82, 255, 255, 255, 37, 8, 17, 192, 0, 134, 0, 0, 197, 119, 255, 58, 251, 47, 1, 130, 0, 5, 1, 72, 0, 0, 255, 66, 255, 66, 65,
		8, 2, 0, 36,
		9, 8, 254, 47, 254, 47, 254, 47, 254, 47,
		5, 1, 65,
	}
	secondSubPacket := genPacket(54, []byte{27, 80}, secondSubDataSegment)

	thirdSubDataSegment := []byte{
		0, 128, 0, 114, 25, 3, 5, 16, 5, 3, 1, 71, 212, 247, 2, 85, 109, 243, 0, 0, 0, 110, 100,
		4, 2, 18, 16,
		2, 44, 68, 1, 255, 255, 255, 67, 251, 3, 249, 183, 251, 3, 102, 68, 255, 255, 255, 37, 8, 17, 173, 3, 178, 2, 184, 197, 75, 254, 218, 249, 11, 2, 222, 255, 90, 1, 72, 0, 58, 255, 2, 255, 56, 65,
		8, 2, 192, 109,
		9, 8, 254, 47, 254, 47, 254, 47, 254, 47,
		5, 1, 67,
	}
	thirdSubPacket := genPacket(54, []byte{27, 81}, thirdSubDataSegment)

	var dataSegment bytes.Buffer
	dataSegment.Write(firstSubPacket)
	dataSegment.Write(secondSubPacket)
	dataSegment.Write(thirdSubPacket)
	packet := genPacket(0x37, []byte{27, 82}, dataSegment.Bytes())
	_, err := writer.Write(packet)
	if err != nil {
		return err
	}
	log.Println("toSendMultiPacketWithAutoReg Write")
	err = writer.Flush()
	if err != nil {
		return err
	}
	log.Println("toSendMultiPacketWithAutoReg Flush")
	return nil
}

func appendLoginItem(dataSegment []byte) []byte {
	var buf bytes.Buffer
	buf.Write(dataSegment)
	buf.WriteByte(1)
	buf.WriteByte(33)
	buf.WriteString(pin)
	buf.WriteString(autoRegDeviceId)
	buf.WriteString(autoRegCntrNum)
	buf.WriteByte('D')

	return buf.Bytes()
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
