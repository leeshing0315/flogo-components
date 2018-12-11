package tcpreceiver

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"strings"
)

const (
	ASCII_CR byte = 13
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
	Type           byte
	Pin            string
	TerminalNum    string
	HardwareVer    string
	Conn           net.Conn
	ServerSocket   *ServerSocket
	RemoteAddrStr  string
	SendCommandSeq uint16
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

	err = handleIpInfo(s, reader)
	if err != nil {
		s.ServerSocket.OnError(s, err)
		s.ServerSocket.OnClose(s)
		s.Conn.Close()
		return
	}

	for {
		binPacket, err := parseByProtocol(s, reader)
		if err != nil {
			s.ServerSocket.OnError(s, err)
			s.ServerSocket.OnClose(s)
			s.Conn.Close()
			return
		}

		err = s.ServerSocket.OnMessage(s, binPacket)
		if err != nil {
			s.ServerSocket.OnError(s, err)
			s.ServerSocket.OnClose(s)
			s.Conn.Close()
			return
		}
	}
}

func handleIpInfo(s *Socket, reader *bufio.Reader) error {
	// PROXY TCP4 <Device IP> <Nginx IP> <Nginx Port> <Local Port> <CR> <LF>
	head, err := reader.ReadBytes(ASCII_CR)
	if err != nil {
		return err
	}
	ipInfo := head[0 : len(head)-1]
	_, err = reader.ReadByte()
	if err != nil {
		return err
	}

	ipInfos := strings.Split(string(ipInfo), " ")
	s.RemoteAddrStr = ipInfos[2]

	return nil
}

func parseByProtocol(s *Socket, reader *bufio.Reader) (*BinPacket, error) {
	command, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	sequence := make([]byte, 2)
	_, err = reader.Read(sequence)
	if err != nil {
		return nil, err
	}

	dataSegmentLength := make([]byte, 2)
	_, err = reader.Read(dataSegmentLength)
	if err != nil {
		return nil, err
	}

	dataSegment := make([]byte, binary.BigEndian.Uint16(dataSegmentLength))
	_, err = reader.Read(dataSegment)
	if err != nil {
		return nil, err
	}

	crc16Check := make([]byte, 2)
	_, err = reader.Read(crc16Check)
	if err != nil {
		return nil, err
	}

	return &BinPacket{
		Command:           command,
		Sequence:          sequence,
		DataSegmentLength: dataSegmentLength,
		DataSegment:       dataSegment,
		CRC16Check:        crc16Check,
	}, nil
}
