package tcpreceiver

import (
	"bufio"
	"encoding/binary"

	"github.com/howeyc/crc16"
)

func handlePacket(socket *Socket, packet *BinPacket) {
	switch packet.Command {
	case 49: // 0x31 - test connection
		handleTest(socket, packet)
	case 50: // 0x32 - login
		handleLogin(socket, packet)
	default:
	}
}

func handleTest(socket *Socket, packet *BinPacket) {
	writer := bufio.NewWriter(socket.Conn)
	content := make([]byte, 7)

	content[0] = 49

	copy(content[1:3], packet.Sequence)

	content[3] = 0
	content[4] = 0

	checksum := crc16.Checksum(content[0:5], crc16.IBMTable)
	binary.BigEndian.PutUint16(content[5:7], checksum)

	writer.Write(content)
	writer.Flush()
}

func handleLogin(socket *Socket, packet *BinPacket) {
	writer := bufio.NewWriter(socket.Conn)
	content := make([]byte, 7)

	content[0] = 50

	copy(content[1:3], packet.Sequence)

	content[3] = 0
	content[4] = 0

	checksum := crc16.Checksum(content[0:5], crc16.IBMTable)
	binary.BigEndian.PutUint16(content[5:7], checksum)

	writer.Write(content)
	writer.Flush()
}
