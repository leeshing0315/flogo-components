package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"log"

	"github.com/sigurn/crc16"
)

func receivePacket(reader *bufio.Reader, writer *bufio.Writer, errChain chan error) {
	for {
		cmd, err := reader.ReadByte()
		if err != nil {
			errChain <- err
			return
		}

		// seqnoBytes := make([]byte, 2)
		// _, err = reader.Read(seqnoBytes)
		seqnoBytes, err := readCount(reader, 2)
		if err != nil {
			errChain <- err
			return
		}

		// dataSegmentLengthBytes := make([]byte, 2)
		// _, err = reader.Read(dataSegmentLengthBytes)
		dataSegmentLengthBytes, err := readCount(reader, 2)
		if err != nil {
			errChain <- err
			return
		}

		dataSegmentLength := binary.BigEndian.Uint16(dataSegmentLengthBytes)
		// dataSegment := make([]byte, dataSegmentLength)
		// n, err := reader.Read(dataSegment)
		dataSegment, err := readCount(reader, int(dataSegmentLength))
		if err != nil {
			errChain <- err
			return
		}

		// crcSegment := make([]byte, 2)
		// _, err = reader.Read(crcSegment)
		crcSegment, err := readCount(reader, 2)
		if err != nil {
			errChain <- err
			return
		}

		switch cmd {
		case 0x34:
			err = receiveSetting(writer, seqnoBytes, dataSegment)
			if err != nil {
				errChain <- err
				return
			}
		case 0x33:
			log.Println("*****receiveFirmware*****")
			log.Println(cmd)
			log.Println(seqnoBytes)
			log.Println(dataSegmentLengthBytes)
			log.Println(dataSegment)
			log.Println(crcSegment)
			log.Println("*****receiveFirmware*****")
			err = receiveFirmware(writer, seqnoBytes, dataSegment)
			if err != nil {
				errChain <- err
				return
			}
		}
	}
}

func receiveSetting(writer *bufio.Writer, seqnoBytes []byte, dataSegment []byte) error {
	if len(dataSegment) == 4 {
		var buf bytes.Buffer
		buf.WriteByte(0x34)
		buf.Write(seqnoBytes)
		buf.WriteByte(0)
		buf.WriteByte(30)
		var dataSegment []byte = make([]byte, 33)
		dataSegment[0] = '*'
		dataSegment[1] = 'L'
		dataSegment[32] = '#'
		buf.Write(dataSegment)
		myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
		checksum := crc16.Checksum(buf.Bytes(), myTable)
		crc := make([]byte, 2)
		binary.LittleEndian.PutUint16(crc, checksum)
		buf.Write(crc)
		_, err := writer.Write(buf.Bytes())
		if err != nil {
			return err
		}
	} else {
		err := responseEmpty(writer, 0x34, seqnoBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

func receiveFirmware(writer *bufio.Writer, seqnoBytes []byte, dataSegment []byte) error {
	seqno := dataSegment[2]
	length := binary.BigEndian.Uint16(dataSegment[3:5])
	content := dataSegment[5 : 5+length]
	switch seqno {
	case 0:
		err := initUpgrade(writer, content, seqnoBytes)
		return err
	case 0xff:
		clearUpgrade()
	default:
		err := receiveFileSlice(seqno, content)
		return err
	}
	return nil
}

func initUpgrade(writer *bufio.Writer, content, seqnoBytes []byte) error {
	log.Println("initUpgrade")
	if len(content) != 118 {
		log.Println("content length not right")
		return errors.New("content length not right")
	}
	firmwareLock.Lock()
	NewFirmware(writer, content, seqnoBytes)
	err := firmwareUpgrade.StartUpgrade()
	firmwareLock.Unlock()
	return err
}

func clearUpgrade() {
	log.Println("clearUpgrade")
	firmwareLock.Lock()
	ClearUpgrade()
	firmwareLock.Unlock()
}

func receiveFileSlice(seqno byte, content []byte) (err error) {
	firmwareLock.Lock()
	defer firmwareLock.Unlock()
	if firmwareUpgrade != nil {
		err = firmwareUpgrade.ReceiveFileSlice(seqno, content)
		if err != nil {
			return err
		}
		if firmwareUpgrade.FileSliceNum == firmwareUpgrade.FileSliceSum {
			firmwareUpgrade.sendEndPacket()
		}
	}
	return err
}

func responseEmpty(writer *bufio.Writer, cmd byte, seqnoBytes []byte) error {
	_, err := writer.Write(genPacket(cmd, seqnoBytes, []byte{}))
	return err
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
