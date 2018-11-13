package tcpreceiver

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
)

// BinPacket the orignal TCP packet
type BinPacket struct {
	Command           byte
	Sequence          []byte
	DataSegmentLength []byte
	DataSegment       []byte
	CRC16Check        byte
}

// Socket TCP client
type Socket struct {
	Pin          string
	TerminalNum  string
	HardwareVer  string
	Conn         net.Conn
	ServerSocket *ServerSocket
}

// ServerSocket TCP server
type ServerSocket struct {
	Address   string
	OnOpen    func(*Socket) error
	OnMessage func(*Socket, *BinPacket) error
	OnClose   func(*Socket)
	OnError   func(*Socket, error)
	Listener  net.Listener
}

// NewServerSocket init a serversocket
func NewServerSocket(address string) *ServerSocket {
	serverSocket := &ServerSocket{
		Address: address,
	}
	serverSocket.OnOpen = func(*Socket) error { return nil }
	serverSocket.OnMessage = func(*Socket, *BinPacket) error { return nil }
	serverSocket.OnClose = func(*Socket) {}
	serverSocket.OnError = func(*Socket, error) {}
	return serverSocket
}

// Listen start to listner port
func (server *ServerSocket) Listen() error {
	listener, err := net.Listen("tcp", server.Address)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		socket := &Socket{
			Conn:         conn,
			ServerSocket: server,
		}
		go socket.execute()
	}
}

func (s *Socket) execute() {
	s.ServerSocket.OnOpen(s)
	reader := bufio.NewReader(s.Conn)
	for {
		command, _ := reader.ReadByte()

		sequence := make([]byte, 2)
		reader.Read(sequence)

		dataSegmentLength := make([]byte, 2)
		reader.Read(dataSegmentLength)

		dataSegment := make([]byte, binary.BigEndian.Uint16(dataSegmentLength))
		s.Conn.Read(dataSegment)

		crc16Check, _ := reader.ReadByte()

		err := s.ServerSocket.OnMessage(s, &BinPacket{
			Command:           command,
			Sequence:          sequence,
			DataSegmentLength: dataSegmentLength,
			DataSegment:       dataSegment,
			CRC16Check:        crc16Check,
		})
		if err != nil {
			s.ServerSocket.OnError(s, err)
		}
	}
}
