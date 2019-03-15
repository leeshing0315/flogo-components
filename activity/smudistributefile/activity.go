package smudistributefile

import (
	"bytes"
	"context"
	"encoding/binary"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/leeshing0315/flogo-components/activity/smuversionchecking"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (a *MyActivity) Eval(ctx activity.Context) (done bool, err error) {

	// do eval
	// firmwareVersionStr := ctx.GetInput("firmwareVersion").(string)
	serialNumber := ctx.GetInput("serialNumber").(int)
	devId := ctx.GetInput("devId").(string)
	uri := ctx.GetInput("uri").(string)
	dbName := ctx.GetInput("dbName").(string)
	identifier := ctx.GetInput("identifier").(string)

	if serialNumber == 0xFF {
		// update firmwareDeployment from inProgress to completed
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		defer client.Disconnect(context.Background())
		if err != nil {
			return false, err
		}
		db := client.Database(dbName)
		coll := db.Collection("firmwareDeployments")
		coll.UpdateOne(
			context.Background(),
			map[string]interface{}{
				"devId":        devId,
				"deployStatus": "inProgress",
			},
			bson.M{
				"$set": map[string]interface{}{
					"deployStatus": "completed",
				},
			},
		)
		ctx.SetOutput("upgradeSegment", generateResponseContent(serialNumber, []byte{}))
		return true, nil
	}

	// firmwareVersion := make(map[string]interface{})
	// json.Unmarshal([]byte(firmwareVersionStr), &firmwareVersion)

	// firmwareFileBytes := getBytesFromMap(firmwareVersion["fileContent"].([]interface{})[0].(map[string]interface{}))
	firmwareFileBytesInterface, ok := smuversionchecking.FirmwareCacheMap.Load(identifier)
	if ok == false {
		return false, nil
	}
	firmwareFileBytes := firmwareFileBytesInterface.([]byte)

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
	firmwareFileBytesLen := len(firmwareFileBytes)
	firstIndex := 512 * (serialNumber - 1)
	lastIndex := 512 * serialNumber
	if firstIndex < firmwareFileBytesLen {
		if lastIndex > firmwareFileBytesLen {
			lastIndex = firmwareFileBytesLen
		}
		copy(firmwareBuff, firmwareFileBytes[firstIndex:lastIndex])
	}

	ctx.SetOutput("upgradeSegment", generateResponseContent(serialNumber, firmwareBuff))
	return true, nil
}

func generateResponseContent(serialNumber int, contentBuff []byte) []byte {
	var upgradeSegmentBuff bytes.Buffer
	upgradeSegmentBuff.WriteString("*Q")
	upgradeSegmentBuff.WriteByte(byte(serialNumber))
	contentLength := make([]byte, 2)
	binary.BigEndian.PutUint16(contentLength, uint16(len(contentBuff)))
	upgradeSegmentBuff.Write(contentLength)
	upgradeSegmentBuff.Write(contentBuff)
	upgradeSegmentBuff.WriteString("#")
	return upgradeSegmentBuff.Bytes()
}

func getBytesFromMap(input map[string]interface{}) []byte {
	resultLen := len(input)
	result := make([]byte, resultLen)
	for k, v := range input {
		index, err := strconv.ParseUint(k, 10, 64)
		if err != nil || int(index) >= resultLen {
			break
		}
		result[index] = byte(v.(float64))
	}
	return result
}
