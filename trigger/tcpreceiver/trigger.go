package tcpreceiver

import (
	"bufio"
	"context"
	"encoding/binary"
	"log"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/sigurn/crc16"
)

// MyTriggerFactory My Trigger factory
type MyTriggerFactory struct {
	metadata *trigger.Metadata
}

// NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &MyTriggerFactory{metadata: md}
}

// New Creates a new trigger instance for a given id
func (t *MyTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &MyTrigger{metadata: t.metadata, config: config}
}

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata     *trigger.Metadata
	config       *trigger.Config
	serverSocket *ServerSocket
	handlers     []*trigger.Handler
}

// Initialize implements trigger.Init.Initialize
func (t *MyTrigger) Initialize(ctx trigger.InitContext) error {
	t.handlers = ctx.GetHandlers()
	return nil
}

// Metadata implements trigger.Trigger.Metadata
func (t *MyTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *MyTrigger) Start() error {
	port := t.config.GetSetting("port")
	url := "0.0.0.0:" + port
	serverSocket := NewServerSocket(url)
	t.serverSocket = serverSocket

	serverSocket.OnOpen = func(s *Socket) error {
		log.Printf("***** Client " + s.Conn.RemoteAddr().String() + " Connected *****")
		return nil
	}

	serverSocket.OnClose = func(s *Socket) {
		log.Printf("***** Client " + s.Conn.RemoteAddr().String() + " Closed *****")
	}

	serverSocket.OnError = func(s *Socket, err error) {
		log.Println(err)
	}

	serverSocket.OnMessage = func(s *Socket, packet *BinPacket) error {
		log.Println(s.Conn.RemoteAddr().String(), packet)

		triggerData := map[string]interface{}{}
		triggerData["eventTime"] = time.Now().Format("2006-01-02 15:04:05.000000")
		triggerData["ip"] = s.Conn.RemoteAddr().String()
		triggerData["command"] = int(packet.Command)
		triggerData["seqNo"] = int(binary.BigEndian.Uint32(packet.Sequence))
		triggerData["reqDataSegment"] = packet.DataSegment
		writer := bufio.NewWriter(s.Conn)
		for _, handler := range t.handlers {
			results, _ := handler.Handle(context.Background(), triggerData)
			if len(results) != 0 {
				dataAttr, ok := results["resDataSegment"]
				if ok {
					dataSegment := dataAttr.Value().([]byte)

					content := make([]byte, len(dataSegment)+7)
					content[0] = packet.Command
					copy(content[1:3], packet.Sequence)
					binary.BigEndian.PutUint16(content[3:5], uint16(len(dataSegment)))
					copy(content[5:5+len(dataSegment)], dataSegment)

					myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
					checksum := crc16.Checksum(content[0:len(content)-2], myTable)
					binary.LittleEndian.PutUint16(content[len(content)-2:len(content)], checksum)

					writer.Write(content)
					writer.Flush()
				}
			}
		}
		return nil
		// err := handlePacket(s, packet)
		// return err
	}
	go t.serverSocket.Listen()
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() error {
	// stop the trigger
	return nil
}
