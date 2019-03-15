package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"

	"github.com/sigurn/crc16"
)

var serverUri string = "localhost:8033"

// var serverUri string = "itciot-tcp.cargosmart.ai:8080"

// var serverUri string = "ciot-tcp.cargosmart.ai:8080"

var deviceIp string = "192.168.0.2"

var localPort string = "8033"

var nginxIp string = "192.168.0.3"

var nginxPort string = "23456"

// var pin string = "460011234567890" // not being registered

var pin string = "460010604706821" // C01937 (15 bytes)

// var pin string = "460011710324088" // C00001 (15 bytes)

var autoRegDeviceId = "C01937" // 6 bytes

// var autoRegDeviceId = "C00001" // 6 bytes

var autoRegCntrNum = "CXRU1338831" // C01937 (11 bytes)

// var autoRegCntrNum = "SMUT0000001" // C00001 (11 bytes)

var firmwareVersion string = "HS19-2"

// var firmwareVersion string = "HS181120"

var defaultElectricalCommunicationFrequency string = "0005"

var defaultNoElectricityCommunicationFrequency string = "0015"

var defaultNoElectricityAcquisitionFrequency string = "0015"

var defaultSettingIpAddress string = "C0A80001" // 192.168.0.1

var defaultSettingIpPort string = "1F61" // 8033

var defualtForcedNotToSleep string = "11"

func main() {

	setDeviceInfo(
		defaultElectricalCommunicationFrequency,
		defaultNoElectricityCommunicationFrequency,
		defaultNoElectricityAcquisitionFrequency,
		defaultSettingIpAddress,
		defaultSettingIpPort,
		defualtForcedNotToSleep,
	)

	conn, err := net.Dial("tcp", serverUri)
	if err != nil {
		return
	}

	bufReader := bufio.NewReader(conn)
	bufWriter := bufio.NewWriter(conn)

	err = sendIpInfo(bufWriter)
	if err != nil {
		return
	}

	loginPacket := genLoginPacket()
	_, err = bufWriter.Write(loginPacket)
	if err != nil {
		return
	}
	err = bufWriter.Flush()
	if err != nil {
		return
	}

	errChain := make(chan error)

	go receivePacket(bufReader, bufWriter, errChain)
	go sendPacket(bufWriter, errChain)

	err = <-errChain
	log.Println(err)
	return
}

func sendIpInfo(writer *bufio.Writer) error {
	ipInfo := genIpInfo()
	_, err := writer.Write(ipInfo)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func genIpInfo() []byte {
	var buf bytes.Buffer
	buf.WriteString("PROXY TCP4 ")
	buf.WriteString(deviceIp)
	buf.WriteString(" ")
	buf.WriteString(nginxIp)
	buf.WriteString(" ")
	buf.WriteString(localPort)
	buf.WriteString(" ")
	buf.WriteString(nginxPort)
	buf.WriteByte(0x0D)
	buf.WriteByte(0x0A)
	return buf.Bytes()
}

func genLoginPacket() []byte {
	var loginDataSegmentBuf bytes.Buffer
	loginDataSegmentBuf.WriteByte(0x02)
	loginDataSegmentBuf.WriteByte(byte(len(pin)))
	loginDataSegmentBuf.WriteString(pin)
	loginDataSegmentBuf.WriteByte(0)
	loginDataSegmentBuf.WriteByte(byte(len(firmwareVersion)))
	loginDataSegmentBuf.WriteString(firmwareVersion)
	loginDataSegment := loginDataSegmentBuf.Bytes()
	var loginPacket bytes.Buffer
	loginPacket.WriteByte(0x32)
	loginPacket.Write([]byte{1, 219})
	loginPacket.WriteByte(0)
	loginPacket.WriteByte(byte(len(loginDataSegment)))
	loginPacket.Write(loginDataSegment)
	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(loginPacket.Bytes(), myTable)
	crc := make([]byte, 2)
	binary.LittleEndian.PutUint16(crc, checksum)
	loginPacket.Write(crc)
	return loginPacket.Bytes()
}
