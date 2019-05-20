package tcpreceiver

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"strconv"
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

	serverSocket.OnMessage = func(s *Socket, writer *bufio.Writer, packet *BinPacket) error {
		log.Println(s.Conn.RemoteAddr().String(), packet)

		triggerData := map[string]interface{}{}
		triggerData["emptyByteArray"] = []byte{}
		loc, _ := time.LoadLocation("Asia/Hong_Kong")
		triggerData["eventTime"] = time.Now().In(loc).Format("2006-01-02T15:04:05+08:00")
		triggerData["ip"] = s.RemoteAddrStr
		triggerData["command"] = int(packet.Command)
		triggerData["seqNo"] = int(binary.BigEndian.Uint16(packet.Sequence))
		triggerData["reqDataSegment"] = packet.DataSegment
		triggerData["cntrNum"] = s.CntrNum
		triggerData["devId"] = s.DevId
		triggerData["pin"] = s.Pin
		triggerData["firmwareVersion"] = s.HardwareVer
		triggerData["devtype"] = s.Type
		triggerData["company"] = s.Company

		var originalPacket bytes.Buffer
		originalPacket.WriteByte(packet.Command)
		originalPacket.Write(packet.Sequence)
		originalPacket.Write(packet.DataSegmentLength)
		originalPacket.Write(packet.DataSegment)
		originalPacket.Write(packet.CRC16Check)
		triggerData["originalPacket"] = originalPacket.Bytes()

		for _, handler := range t.handlers {
			results, err := handler.Handle(context.Background(), triggerData)
			if err != nil {
				return err
			}
			if len(results) != 0 {
				dataAttr, ok := results["resDataSegment"]
				setCommandAttr, _ := results["setCommandSegment"]
				readCommandAttr, _ := results["readCommandSegment"]
				upgradeAttr, _ := results["upgradeSegment"]
				cntrNumAttr, _ := results["cntrNum"]
				devIdAttr, _ := results["devId"]
				pinAttr, _ := results["pin"]
				firmwareVersionAttr, _ := results["firmwareVersion"]
				devtypeAttr, _ := results["devtype"]
				stopUpgradeSegmentAttr, _ := results["stopUpgradeSegment"]
				companyAttr, _ := results["company"]

				cntrNum := cntrNumAttr.Value().(string)
				if cntrNum != "" {
					s.CntrNum = cntrNum
				}
				devId := devIdAttr.Value().(string)
				if devId != "" {
					s.DevId = devId
				}
				pin := pinAttr.Value().(string)
				if pin != "" {
					s.Pin = pin
				}
				firmwareVersion := firmwareVersionAttr.Value().(string)
				if firmwareVersion != "" {
					s.HardwareVer = firmwareVersion
				}
				devtype := devtypeAttr.Value().(string)
				if devtype != "" {
					devtypeStr, _ := strconv.ParseInt(devtype, 10, 64)
					s.Type = byte(devtypeStr)
				}
				company := companyAttr.Value().(string)
				if company != "" {
					s.Company = company
				}

				if ok && (packet.Command != 0x34 && !(packet.Command == 0x33 && packet.DataSegment[1] == 'L')) {
					dataSegment := dataAttr.Value().([]byte)
					err := writeToDevice(packet, writer, dataSegment)
					if err != nil {
						return err
					}
				}

				if setCommandAttr.Value() != nil {
					setCommand := setCommandAttr.Value().([]byte)
					if len(setCommand) != 0 {
						commandSeqAttr, _ := results["setCommandSeqNo"]
						commandSeqNo := commandSeqAttr.Value().(string)
						commandSeqNoUint, _ := strconv.ParseUint(commandSeqNo, 10, 16)
						err := sendCommandToDevice(0x34, uint16(commandSeqNoUint), writer, setCommand, false)
						if err != nil {
							return err
						}
					}
				}

				if stopUpgradeSegmentAttr.Value() != nil {
					stopUpgradeSegment := stopUpgradeSegmentAttr.Value().([]byte)
					if len(stopUpgradeSegment) != 0 && (setCommandAttr.Value() == nil || len(setCommandAttr.Value().([]byte)) == 0) {
						err := sendCommandToDevice(0x34, uint16(0xFA), writer, stopUpgradeSegment, false)
						if err != nil {
							return err
						}
					}
				}

				if readCommandAttr.Value() != nil {
					readCommand := readCommandAttr.Value().([]byte)
					if len(readCommand) != 0 {
						commandSeqAttr, _ := results["readCommandSeqNo"]
						commandSeqNo := commandSeqAttr.Value().(string)
						commandSeqNoUint, _ := strconv.ParseUint(commandSeqNo, 10, 16)
						err := sendCommandToDevice(0x34, uint16(commandSeqNoUint), writer, readCommand, false)
						if err != nil {
							return err
						}
					}
				}

				if upgradeAttr.Value() != nil {
					upgradeCommand := upgradeAttr.Value().([]byte)
					if len(upgradeCommand) != 0 {
						// commandSeqAttr, _ := results["upgradeSeqNo"]
						// commandSeqNo := commandSeqAttr.Value().(string)
						// commandSeqNoUint, _ := strconv.ParseUint(commandSeqNo, 10, 16)
						// err := sendCommandToDevice(0x34, uint16(commandSeqNoUint), writer, upgradeCommand)
						err := sendCommandToDevice(0x34, s.CommandSeq, writer, upgradeCommand, true)
						s.CommandSeq = s.CommandSeq + 1
						if err != nil {
							return err
						}
					}
				}
			}
		}
		return nil
	}
	go t.serverSocket.Listen()
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() error {
	// stop the trigger
	return nil
}

func writeToDevice(packet *BinPacket, writer *bufio.Writer, dataSegment []byte) error {
	content := make([]byte, len(dataSegment)+7)
	if packet.Command == 0x33 {
		content[0] = 0x34
	} else {
		content[0] = packet.Command
	}
	copy(content[1:3], packet.Sequence)
	binary.BigEndian.PutUint16(content[3:5], uint16(len(dataSegment)))
	copy(content[5:5+len(dataSegment)], dataSegment)

	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(content[0:len(content)-2], myTable)
	if packet.Command == 0x33 {
		binary.BigEndian.PutUint16(content[len(content)-2:len(content)], checksum)
	} else {
		binary.LittleEndian.PutUint16(content[len(content)-2:len(content)], checksum)
	}
	log.Println("**********Ack:", convertBytesToStrings(content), "**********")

	_, err := writer.Write(content)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func convertBytesToStrings(input []byte) []string {
	output := make([]string, len(input))
	for index, val := range input {
		output[index] = "0x" + strconv.FormatUint(uint64(val), 16)
	}
	return output
}

func sendCommandToDevice(cmdValue int, seqNo uint16, writer *bufio.Writer, dataSegment []byte, crcRotate bool) error {
	content := make([]byte, len(dataSegment)+7)
	content[0] = byte(cmdValue)
	binary.BigEndian.PutUint16(content[1:3], seqNo)
	binary.BigEndian.PutUint16(content[3:5], uint16(len(dataSegment)))
	copy(content[5:5+len(dataSegment)], dataSegment)

	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(content[0:len(content)-2], myTable)
	if crcRotate == true {
		binary.BigEndian.PutUint16(content[len(content)-2:len(content)], checksum)
	} else {
		binary.LittleEndian.PutUint16(content[len(content)-2:len(content)], checksum)
	}

	log.Println("**********Cmd:", convertBytesToStrings(content), "**********")

	_, err := writer.Write(content)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
