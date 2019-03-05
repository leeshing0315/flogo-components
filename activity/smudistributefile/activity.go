package smudistributefile

import (
	"bytes"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval
	firmwareVersionStr := context.GetInput("firmwareVersion").(string)
	serialNumber := context.GetInput("serialNumber").(int)

	if serialNumber == 0xFF {
		context.SetOutput("upgradeSegment", generateResponseContent([]byte{}))
		return true, nil
	}

	firmwareVersion := make(map[string]interface{})
	json.Unmarshal([]byte(firmwareVersionStr), &firmwareVersion)

	firmwareFileBytes := firmwareVersion["firmwareFile"].([]byte)

	// filePath := firmwareVersion["filePath"]
	// // Get file
	// file, err := os.Open(filePath)
	// if err != nil {
	// 	return false, err
	// }
	// defer file.Close()

	// bufReader := bufio.NewReader(file)

	// bufReader.Discard(512 * (serialNumber - 1))
	// firmwareBuff := make([]byte, 512)
	// _, err = bufReader.Read(firmwareBuff)
	// if err != nil && err != io.EOF {
	// 	return false, err
	// }

	firmwareBuff := make([]byte, 512)
	copy(firmwareBuff, firmwareFileBytes[512*(serialNumber-1):512*serialNumber])

	context.SetOutput("upgradeSegment", generateResponseContent(firmwareBuff))
	return true, nil
}

func generateResponseContent(contentBuff []byte) []byte {
	var upgradeSegmentBuff bytes.Buffer
	upgradeSegmentBuff.WriteString("*Q")
	upgradeSegmentBuff.Write(contentBuff)
	upgradeSegmentBuff.WriteString("#")
	return upgradeSegmentBuff.Bytes()
}
