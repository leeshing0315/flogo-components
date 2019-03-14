package smudistributefile

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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
	firmwareVersionStr := ctx.GetInput("firmwareVersion").(string)
	serialNumber := ctx.GetInput("serialNumber").(int)
	devId := ctx.GetInput("devId").(string)
	uri := ctx.GetInput("uri").(string)
	dbName := ctx.GetInput("dbName").(string)

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
				"devId":  devId,
				"status": "inProgress",
			},
			bson.M{
				"$set": map[string]interface{}{
					"status": "completed",
				},
			},
		)
		ctx.SetOutput("upgradeSegment", generateResponseContent([]byte{}))
		return true, nil
	}

	firmwareVersion := make(map[string]interface{})
	json.Unmarshal([]byte(firmwareVersionStr), &firmwareVersion)

	firmwareFileBytes := getBytesFromMap(firmwareVersion["fileContent"].([]interface{})[0].(map[string]interface{}))

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

	ctx.SetOutput("upgradeSegment", generateResponseContent(firmwareBuff))
	return true, nil
}

func generateResponseContent(contentBuff []byte) []byte {
	var upgradeSegmentBuff bytes.Buffer
	upgradeSegmentBuff.WriteString("*Q")
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
		result[index] = byte(v.(int32))
	}
	return result
}
