package main

import (
	"bufio"
	"bytes"
	"encoding/binary"

	"github.com/sigurn/crc16"
)

func receivePacket(reader *bufio.Reader, writer *bufio.Writer, errChain chan error) {
	for {
		cmd, err := reader.ReadByte()
		if err != nil {
			errChain <- err
			return
		}

		seqnoBytes := make([]byte, 2)
		_, err = reader.Read(seqnoBytes)
		if err != nil {
			errChain <- err
			return
		}

		dataSegmentLengthBytes := make([]byte, 2)
		_, err = reader.Read(dataSegmentLengthBytes)
		if err != nil {
			errChain <- err
			return
		}

		dataSegmentLength := binary.BigEndian.Uint16(dataSegmentLengthBytes)
		dataSegment := make([]byte, dataSegmentLength)
		_, err = reader.Read(dataSegment)
		if err != nil {
			errChain <- err
			return
		}

		crcSegment := make([]byte, 2)
		_, err = reader.Read(crcSegment)
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
	return nil
}

func responseEmpty(writer *bufio.Writer, cmd byte, seqnoBytes []byte) error {
	_, err := writer.Write(genPacket(cmd, seqnoBytes, []byte{}))
	return err
}
