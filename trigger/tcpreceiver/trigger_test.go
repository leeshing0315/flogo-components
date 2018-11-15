package tcpreceiver

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/stretchr/testify/mock"
)

func getJsonMetadata() string {
	jsonMetadataBytes, err := ioutil.ReadFile("trigger.json")
	if err != nil {
		panic("No Json Metadata found for trigger.json path")
	}
	return string(jsonMetadataBytes)
}

const testConfig string = `{
  "id": "mytrigger",
  "settings": {
    "port": "8033"
  },
  "handlers": [
    {
      "settings": {
        "handler_setting": "somevalue"
      },
      "action" {
	     "id": "test_action"
      }
    }
  ]
}`

func TestCreate(t *testing.T) {

	// New factory
	md := trigger.NewMetadata(getJsonMetadata())
	f := NewFactory(md)

	if f == nil {
		t.Fail()
	}

	// New Trigger
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)
	trg := f.New(&config)
	trg.Start()

	if trg == nil {
		t.Fail()
	}
}

type MockedConn struct {
	mock.Mock
}

func (m *MockedConn) Read(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockedConn) Write(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockedConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedConn) LocalAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *MockedConn) RemoteAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *MockedConn) SetDeadline(t time.Time) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedConn) SetReadDeadline(t time.Time) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedConn) SetWriteDeadline(t time.Time) error {
	args := m.Called()
	return args.Error(0)
}

func TestHandlerTest(t *testing.T) {
	mockedConn := new(MockedConn)
	mockedConn.On("Write", mock.Anything).Return(1, nil)

	socket := &Socket{
		Conn: mockedConn,
	}

	packet := &BinPacket{
		Command:           49,
		Sequence:          []byte{0x00, 0x07},
		DataSegmentLength: []byte{0x00, 0x00},
		DataSegment:       []byte{},
		CRC16Check:        []byte{0x00, 0x01},
	}

	handlePacket(socket, packet)

	// fmt.Println(mockedConn.AssertCalled(t, "write", mock.Anything))
}
