package smudistributefile

import (
	"bytes"
	"context"
	"encoding/binary"
	"strconv"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/leeshing0315/flogo-components/common/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// store firmware for cache
var firmwareCacheMap = make(map[string][]byte)
var firmwareCacheLock sync.RWMutex

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(ctx activity.Context) (done bool, err error) {

	// do eval
	reqDataSegmentBytes := ctx.GetInput("reqDataSegment").([]byte)
	devId := ctx.GetInput("devId").(string)
	uri := ctx.GetInput("uri").(string)
	dbName := ctx.GetInput("dbName").(string)

	if reqDataSegmentBytes[1] == 'L' {
		// update read config content
		value := service.DecodeReadConfigAck(reqDataSegmentBytes)
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		defer client.Disconnect(context.Background())
		if err != nil {
			return false, err
		}
		db := client.Database(dbName)
		cmdColl := db.Collection("remotecommands")
		historyCmdColl := db.Collection("remotecommandhistories")
		cmdColl.FindOneAndUpdate(
			context.Background(),
			bson.M{
				"devid":    devId,
				"subcmd":   "FF",
				"sendflag": "2",
			},
			bson.M{
				"$set": map[string]interface{}{
					"value": value,
				},
			},
			&options.FindOneAndUpdateOptions{
				Sort: bson.M{
					"savetime": -1,
				},
			},
		)
		historyCmdColl.FindOneAndUpdate(
			context.Background(),
			bson.M{
				"devid":    devId,
				"subcmd":   "FF",
				"sendflag": "2",
			},
			bson.M{
				"$set": map[string]interface{}{
					"value": value,
				},
			},
			&options.FindOneAndUpdateOptions{
				Sort: bson.M{
					"savetime": -1,
				},
			},
		)
		ctx.SetOutput("upgradeSegment", []byte{})
		return true, nil
	}

	serialNumber := int(reqDataSegmentBytes[2])
	identifier := string(reqDataSegmentBytes[5:13])

	if serialNumber == 0xFF {
		// update firmwareDeployment from inProgress to completed
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		defer client.Disconnect(context.Background())
		if err != nil {
			return false, err
		}
		db := client.Database(dbName)
		coll := db.Collection("firmwareDeployments")
		coll.UpdateMany(
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

	var firmwareFileBytes []byte
	firmwareCacheLock.RLock()
	firmwareFileBytes = firmwareCacheMap[identifier]
	firmwareCacheLock.RUnlock()
	if len(firmwareFileBytes) == 0 {
		firmwareVersion := make(map[string]interface{})
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		defer client.Disconnect(context.Background())
		if err != nil {
			return false, err
		}
		db := client.Database(dbName)
		coll := db.Collection("firmwareVersions")
		singleResult := coll.FindOne(context.Background(), bson.M{"identifier": identifier})
		err = singleResult.Decode(&firmwareVersion)
		if err != nil {
			return false, nil
		}
		firmwareFileBytes = getBytesFromMap(firmwareVersion["fileContent"].(primitive.A)[0].(map[string]interface{}))
		firmwareCacheLock.Lock()
		firmwareCacheMap[identifier] = firmwareFileBytes
		firmwareCacheLock.Unlock()
	}

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
		result[index] = byte(v.(int32))
	}
	return result
}
