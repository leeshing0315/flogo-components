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
	CRC16Check        []byte
}

// Socket TCP client
type Socket struct {
	Type         byte
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
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		socket := &Socket{
			Conn:         conn,
			ServerSocket: server,
		}
		go socket.execute()
	}
}

func (s *Socket) execute() {
	err := s.ServerSocket.OnOpen(s)
	if err != nil {
		s.ServerSocket.OnError(s, err)
		s.ServerSocket.OnClose(s)
		s.Conn.Close()
		return
	}
	reader := bufio.NewReader(s.Conn)
	for {
		command, err := reader.ReadByte()
		if err != nil {
			s.ServerSocket.OnError(s, err)
			s.ServerSocket.OnClose(s)
			s.Conn.Close()
			return
		}

		sequence := make([]byte, 2)
		_, err = reader.Read(sequence)
		if err != nil {
			s.ServerSocket.OnError(s, err)
			s.ServerSocket.OnClose(s)
			s.Conn.Close()
			return
		}

		dataSegmentLength := make([]byte, 2)
		_, err = reader.Read(dataSegmentLength)
		if err != nil {
			s.ServerSocket.OnError(s, err)
			s.ServerSocket.OnClose(s)
			s.Conn.Close()
			return
		}

		dataSegment := make([]byte, binary.BigEndian.Uint16(dataSegmentLength))
		_, err = s.Conn.Read(dataSegment)
		if err != nil {
			s.ServerSocket.OnError(s, err)
			s.ServerSocket.OnClose(s)
			s.Conn.Close()
			return
		}

		crc16Check := make([]byte, 2)
		_, err = reader.Read(crc16Check)
		if err != nil {
			s.ServerSocket.OnError(s, err)
			s.ServerSocket.OnClose(s)
			s.Conn.Close()
			return
		}

		err = s.ServerSocket.OnMessage(s, &BinPacket{
			Command:           command,
			Sequence:          sequence,
			DataSegmentLength: dataSegmentLength,
			DataSegment:       dataSegment,
			CRC16Check:        crc16Check,
		})
		if err != nil {
			s.ServerSocket.OnError(s, err)
			s.ServerSocket.OnClose(s)
			s.Conn.Close()
			return
		}
	}
}
