package tcpreceiver

import (
	"log"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
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
}

// Initialize implements trigger.Init.Initialize
func (t *MyTrigger) Initialize(ctx trigger.InitContext) error {
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
		err := handlePacket(s, packet)
		return err
	}

	err := serverSocket.Listen()

	return err
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() error {
	// stop the trigger
	return nil
}
