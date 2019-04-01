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
	Type          byte
	Pin           string
	TerminalNum   string
	HardwareVer   string
	Conn          net.Conn
	ServerSocket  *ServerSocket
	RemoteAddrStr string
	CntrNum       string
	DevId         string
	CommandSeq    uint16
}

// ServerSocket TCP server
type ServerSocket struct {
	Address   string
	OnOpen    func(*Socket) error
	OnMessage func(*Socket, *bufio.Writer, *BinPacket) error
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
	serverSocket.OnMessage = func(*Socket, *bufio.Writer, *BinPacket) error { return nil }
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
			CommandSeq:   1,
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
	writer := bufio.NewWriter(s.Conn)

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

		err = s.ServerSocket.OnMessage(s, writer, binPacket)
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

	sequence, err := readCount(reader, 2)
	if err != nil {
		return nil, err
	}

	dataSegmentLength, err := readCount(reader, 2)
	if err != nil {
		return nil, err
	}

	dataSegment, err := readCount(reader, int(binary.BigEndian.Uint16(dataSegmentLength)))
	if err != nil {
		return nil, err
	}

	crc16Check, err := readCount(reader, 2)
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
