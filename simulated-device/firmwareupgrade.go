package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"sync"
	"time"
)

type FirmwareUpgrade struct {
	Identifier      string
	FileLen         uint32
	FileCrc         uint16
	FileName        string
	UpgradeDate     string
	UpgradeTime     string
	FirmwareVersion string
	Operator        string
	FileDescription string
	FileSlices      []*FileSlice
	FileSliceSum    byte
	FileSliceNum    byte
	Writer          *bufio.Writer
	SeqnoBytes      []byte
}

type FileSlice struct {
	Identifier  string
	Seqno       byte
	FirmwareObj []byte
}

var firmwareLock sync.Mutex

var firmwareUpgrade *FirmwareUpgrade

func NewFirmware(writer *bufio.Writer, content, seqnoBytes []byte) {
	if firmwareUpgrade != nil {
		firmwareUpgrade = nil
	}
	firmwareUpgrade = &FirmwareUpgrade{
		Identifier:      string(content[0:8]),
		FileLen:         binary.BigEndian.Uint32(content[8:12]),
		FileCrc:         binary.BigEndian.Uint16(content[12:14]),
		FileName:        string(content[14:22]),
		UpgradeDate:     string(content[22:30]),
		UpgradeTime:     string(content[30:38]),
		FirmwareVersion: string(content[38:46]),
		Operator:        string(content[46:54]),
		FileDescription: string(content[54:118]),
		FileSlices:      make([]*FileSlice, binary.BigEndian.Uint32(content[8:12])/512+1),
		FileSliceSum:    byte(binary.BigEndian.Uint32(content[8:12])/512 + 1),
		FileSliceNum:    0,
		Writer:          writer,
		SeqnoBytes:      seqnoBytes,
	}
}

func (fu *FirmwareUpgrade) StartUpgrade() error {
	for i := 1; i <= int(fu.FileSliceSum); i++ {
		err := fu.requestFirmware(byte(i))
		if err != nil {
			return err
		}
	}
	go fu.goCheckUpgradeProcess()
	return nil
}

func (fu *FirmwareUpgrade) goCheckUpgradeProcess() {
	startTime := time.Now().Unix()
	const MAX_UPGRADE_SECONDS = 300
	reSendTime := startTime
	const MAX_RESEND_SECONDS = 60
	for {
		nowTime := time.Now().Unix()
		if nowTime-reSendTime > MAX_UPGRADE_SECONDS {
			// 超时结束
			if fu.FileSliceNum == fu.FileSliceSum {
				fu.sendEndPacket()
				reSendTime = nowTime
				break
			}
			firmwareLock.Lock()
			ClearUpgrade()
			defer firmwareLock.Unlock()
			return
		}
		if nowTime-reSendTime > MAX_RESEND_SECONDS {
			// 请求重发
			for index, slice := range fu.FileSlices {
				if slice != nil {
					break
				}
				err := fu.requestFirmware(byte(index))
				if err != nil {
					return
				}
			}
			reSendTime = nowTime
		}
	}
}

func (fu *FirmwareUpgrade) sendEndPacket() error {
	var buf bytes.Buffer
	buf.WriteString("*Q")
	buf.WriteByte(0xFF)
	buf.WriteByte(0x00)
	buf.WriteByte(0x08)
	buf.WriteString(fu.Identifier)
	buf.WriteString("#")
	_, err := fu.Writer.Write(genPacket(0x33, fu.SeqnoBytes, buf.Bytes()))
	if err != nil {
		return err
	}
	err = fu.Writer.Flush()
	return err
}

func (fu *FirmwareUpgrade) requestFirmware(seqno byte) error {
	var buf bytes.Buffer
	buf.WriteString("*Q")
	buf.WriteByte(seqno)
	buf.WriteByte(0x00)
	buf.WriteByte(0x08)
	buf.WriteString(fu.Identifier)
	buf.WriteString("#")
	log.Println("before send")
	result := genPacket(0x33, fu.SeqnoBytes, buf.Bytes())
	log.Println(result)
	_, err := fu.Writer.Write(result)
	if err != nil {
		return err
	}
	err = fu.Writer.Flush()
	log.Println("after send")
	return err
}

func (fu *FirmwareUpgrade) ReceiveFileSlice(seqno byte, content []byte) error {
	if fu.FileSlices[seqno-1] == nil {
		fu.FileSlices[seqno-1] = &FileSlice{
			Identifier:  fu.Identifier,
			Seqno:       seqno,
			FirmwareObj: content,
		}
		fu.FileSliceNum++
	}
	return nil
}

func ClearUpgrade() {
	firmwareUpgrade = nil
}
