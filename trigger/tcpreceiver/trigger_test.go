package tcpreceiver

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/sigurn/crc16"
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

type MockedAddr struct {
	mock.Mock
}

func (m *MockedAddr) Network() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockedAddr) String() string {
	args := m.Called()
	return args.String(0)
}

func TestHandlerTest(t *testing.T) {
	mockedConn := new(MockedConn)
	mockedAddr := new(MockedAddr)
	mockedAddr.On("Network").Return("tcp")
	mockedAddr.On("String").Return("127.0.0.1")
	mockedConn.On("RemoteAddr").Return(mockedAddr)
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

func TestHandlerLogin(t *testing.T) {
	mockedConn := new(MockedConn)
	mockedAddr := new(MockedAddr)
	mockedAddr.On("Network").Return("tcp")
	mockedAddr.On("String").Return("127.0.0.1")
	mockedConn.On("RemoteAddr").Return(mockedAddr)
	mockedConn.On("Write", mock.Anything).Return(1, nil)

	socket := &Socket{
		Conn: mockedConn,
	}

	packet := &BinPacket{
		Command:           50,
		Sequence:          []byte{12, 207},
		DataSegmentLength: []byte{0, 27},
		DataSegment:       []byte{2, 15, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56, 0, 8, 72, 83, 49, 56, 49, 49, 48, 56},
		CRC16Check:        []byte{17, 112},
	}

	handlePacket(socket, packet)

	// fmt.Println(mockedConn.AssertCalled(t, "write", mock.Anything))
}

func TestCRC16_1(t *testing.T) {
	packet := []byte{50, 12, 207, 0, 27, 2, 15, 52, 54, 48, 48, 49, 49, 55, 49, 48, 51, 50, 52, 48, 56, 56, 0, 8, 72, 83, 49, 56, 49, 49, 48, 56}
	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(packet, myTable)
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, checksum)
	println(result)
}

func TestCRC16_2(t *testing.T) {
	packet := []byte{49, 0, 7, 0, 0}
	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(packet, myTable)
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, checksum)
	println(result)
}
