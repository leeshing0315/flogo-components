package smuversionchecking

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/sigurn/crc16"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata         *activity.Metadata
	mongoClient      *mongo.Client
	clientGetterLock sync.Mutex
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
	uri := ctx.GetInput("uri").(string)
	dbName := ctx.GetInput("dbName").(string)
	deviceId := ctx.GetInput("devId").(string)

	if deviceId == "" {
		return true, nil
	}

	db, err := getDBConnection(a, uri, dbName)

	if err != nil {
		log.Printf("DB Connect fail: %v", err.Error())
		return false, err
	}
	// do eval
	err, deploymentMap, firmwareVersionMap := queryDeviceFirmwareInformation(db, deviceId)
	if err != nil {
		return false, err
	}
	if len(deploymentMap) == 0 || len(firmwareVersionMap) == 0 || !isReachedDeploymentDate(deploymentMap["targetDeployDate"].(string)) {
		return true, nil
	}
	// Response upgrade command and update upgrade status
	upgradeCommand := handleUpgradeCommand(firmwareVersionMap)

	err = updateDeviceDeploymentStatus(db, deviceId)
	if err != nil {
		log.Printf("Update deployment status fail: %v", err.Error())
		return false, err
	}

	var upgradeSegmentBuff bytes.Buffer
	upgradeSegmentBuff.WriteString("*Q")
	upgradeSegmentBuff.Write(upgradeCommand)
	upgradeSegmentBuff.WriteString("#")

	ctx.SetOutput("upgradeSegment", upgradeSegmentBuff.Bytes())
	return true, nil
}

func getDBConnection(a *MyActivity, uri string, dbName string) (*mongo.Database, error) {
	client := a.mongoClient
	var err error
	if client == nil {
		a.clientGetterLock.Lock()
		client = a.mongoClient
		if client == nil {
			client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
			if err != nil {
				log.Printf("Connection error: %v", err)
				a.clientGetterLock.Unlock()
				return nil, err
			}
			a.mongoClient = client
		}
		a.clientGetterLock.Unlock()
	}

	db := client.Database(dbName)
	return db, nil
}

func queryDeviceFirmwareInformation(db *mongo.Database, devId string) (error, map[string]interface{}, map[string]interface{}) {
	deploymentsColl := db.Collection("firmwareDeployments")

	// if find a firmwareDeployment is inProgress, don't start a new upgrade
	// dpInprogressFilter := buildDeploymentInProgressBsonFilter(devId)
	// deploymentInprogressMap := make(map[string]interface{})
	// err := deploymentsColl.FindOne(context.Background(), dpInprogressFilter).Decode(&deploymentInprogressMap)
	// if err == nil {
	// 	return nil, nil, nil
	// }

	dpBsonFilter := buildBsonFilter(devId)
	deploymentMap := make(map[string]interface{})
	err := deploymentsColl.FindOne(context.Background(), dpBsonFilter).Decode(&deploymentMap)
	if err != nil {
		log.Printf("Connection query firmware error: %v", err)
		return nil, nil, nil
	}

	firmwareVersion := deploymentMap["firmwareVersion"].(string)
	firmwareVersionColl := db.Collection("firmwareVersions")
	fvBsonFilter := buildFirmwareVersionFilter(firmwareVersion)
	firmwareVersionMap := make(map[string]interface{})
	err = firmwareVersionColl.FindOne(context.Background(), fvBsonFilter).Decode(&firmwareVersionMap)
	if err != nil {
		log.Printf("Connection query firmware versions error: %v", err)
		return nil, nil, nil
	}

	return nil, deploymentMap, firmwareVersionMap
}

func isReachedDeploymentDate(targetDeployDateStr string) bool {
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	currentDateStr := time.Now().In(loc).Format("2006-01-02T15:04:05+08:00")
	return strings.Compare(currentDateStr, targetDeployDateStr) >= 0
}

func buildFirmwareVersionFilter(firmwareVersion string) bson.M {
	return bson.M{"firmwareVersion": firmwareVersion}
}

func buildBsonFilter(devId string) bson.M {
	return bson.M{"devId": devId, "deployStatus": "pending"}
}

func buildDeploymentInProgressBsonFilter(devId string) bson.M {
	return bson.M{"devId": devId, "deployStatus": "inProgress"}
}

// store firmware for cache
// var FirmwareCacheMap sync.Map

func handleUpgradeCommand(firmwareVersionMap map[string]interface{}) []byte {
	// firmwareName := []byte(firmwareVersionMap["firmwareName"])
	firmwareLength, _ := strconv.Atoi(firmwareVersionMap["fileSize"].(string))
	// crcValue, _ := strconv.Atoi(firmwareVersionMap["crcValue"].(string))
	firmwareVersion := []byte(firmwareVersionMap["firmwareVersion"].(string))
	operator := []byte(firmwareVersionMap["createdBy"].(string))
	firmwareDesc := []byte(firmwareVersionMap["releaseNote"].(string))
	identifier := []byte(firmwareVersionMap["identifier"].(string))

	// cal crcValue by myself
	fileContentBytes := getBytesFromMap(firmwareVersionMap["fileContent"].(primitive.A)[0].(map[string]interface{}))
	// store firmware for cache
	// FirmwareCacheMap.Store(string(identifier), fileContentBytes)
	myTable := crc16.MakeTable(crc16.CRC16_MODBUS)
	checksum := crc16.Checksum(fileContentBytes, myTable)
	crcValue := int(checksum)

	serialNumber := []byte{0x00}

	contentLength := make([]byte, 2)
	binary.BigEndian.PutUint16(contentLength, 118)

	// download firmware length
	downloadFirmwareLength := make([]byte, 4)
	binary.BigEndian.PutUint32(downloadFirmwareLength, uint32(firmwareLength))

	// crcValue
	crcValueTag := make([]byte, 2)
	binary.LittleEndian.PutUint16(crcValueTag, uint16(crcValue))

	// firmware name
	firmwareNameTag := []byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}
	if len(firmwareVersion) < 8 {
		copy(firmwareNameTag[0:len(firmwareVersion)], firmwareVersion)
	} else {
		copy(firmwareNameTag[0:8], firmwareVersion)
	}

	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	// upgrade date
	upgradeDate := make([]byte, 8)
	copy(upgradeDate, time.Now().In(loc).Format("06-01-02"))
	// upgrade time
	upgradeTime := make([]byte, 8)
	copy(upgradeTime, time.Now().In(loc).Format("03:04:05"))

	// software version
	firmwareVersionTag := make([]byte, 8)
	copy(firmwareVersionTag[0:8], firmwareNameTag)

	// operator
	operatorTag := []byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}
	if len(operator) < 8 {
		copy(operatorTag[0:len(operator)], operator)
	} else {
		copy(operatorTag[0:8], operator)
	}

	// file desc
	const totalFileLength = 64
	firmwareDescTag := make([]byte, 64)
	for i := range firmwareDescTag {
		firmwareDescTag[i] = 0x20
	}
	if len(firmwareDesc) < totalFileLength {
		copy(firmwareDescTag[0:len(firmwareDesc)], firmwareDesc)
	} else {
		copy(firmwareDescTag[0:totalFileLength], firmwareDesc)
	}

	var contentBuff bytes.Buffer
	contentBuff.Write(serialNumber)
	contentBuff.Write(contentLength)
	contentBuff.Write(identifier)
	contentBuff.Write(downloadFirmwareLength)
	contentBuff.Write(crcValueTag)
	contentBuff.Write(firmwareNameTag)
	contentBuff.Write(upgradeDate)
	contentBuff.Write(upgradeTime)
	contentBuff.Write(firmwareVersionTag)
	contentBuff.Write(operatorTag)
	contentBuff.Write(firmwareDescTag)

	return contentBuff.Bytes()
}

func updateDeviceDeploymentStatus(db *mongo.Database, deviceId string) error {
	valueMap := make(map[string]string)
	valueMap["deployStatus"] = "inProgress"
	valueMap["targetDeployDate"] = time.Now().Format("2006-01-02 15:04:05")

	firmwareDeploymentColl := db.Collection("firmwareDeployments")
	_, err := firmwareDeploymentColl.UpdateOne(
		context.Background(),
		bson.M{
			"devId":        deviceId,
			"deployStatus": "pending",
		},
		bson.M{
			"$set": valueMap,
		},
	)
	if err != nil {
		return err
	}
	return nil
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
