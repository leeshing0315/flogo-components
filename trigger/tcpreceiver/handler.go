package tcpreceiver

import (
	"bufio"
	"encoding/binary"
	"errors"
	"log"

	"github.com/sigurn/crc16"
)

func handlePacket(socket *Socket, packet *BinPacket) error {
	switch packet.Command {
	case 49: // 0x31 - test connection
		return handleTest(socket, packet)
	case 50: // 0x32 - login
		return handleLogin(socket, packet)
	default:
		return errors.New("Command mismatch")
	}
}

func handleTest(socket *Socket, packet *BinPacket) error {
	log.Println("***** Client " + socket.Conn.RemoteAddr().String() + " send Command 0x31 *****")

	writer := bufio.NewWriter(socket.Conn)
	content := make([]byte, 7)

	content[0] = 49 // 0x31

	copy(content[1:3], packet.Sequence)

	content[3] = 0
	content[4] = 0

	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(content[0:5], myTable)
	binary.LittleEndian.PutUint16(content[5:7], checksum)

	writer.Write(content)
	writer.Flush()
	return nil
}

func handleLogin(socket *Socket, packet *BinPacket) error {
	log.Println("***** Client " + socket.Conn.RemoteAddr().String() + " send Command 0x32 *****")

	var data = packet.DataSegment

	var cursor int
	socket.Type = data[cursor]
	cursor++
	pinLen := int(data[cursor : cursor+1][0])
	cursor++
	socket.Pin = string(data[cursor : cursor+pinLen])
	cursor += pinLen
	terminalNumLen := int(data[cursor : cursor+1][0])
	cursor++
	socket.TerminalNum = string(data[cursor : cursor+terminalNumLen])
	cursor += terminalNumLen
	hardwareVerLen := int(data[cursor : cursor+1][0])
	cursor++
	socket.HardwareVer = string(data[cursor : cursor+hardwareVerLen])

	writer := bufio.NewWriter(socket.Conn)
	content := make([]byte, 7)

	content[0] = 50 // 0x32

	copy(content[1:3], packet.Sequence)

	content[3] = 0
	content[4] = 0

	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(content[0:5], myTable)
	binary.LittleEndian.PutUint16(content[5:7], checksum)

	writer.Write(content)
	writer.Flush()
	return nil
}
