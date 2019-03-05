package smuversionchecking

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
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
	if len(deploymentMap) == 0 || len(firmwareVersionMap) == 0 || !isReachedDeploymentDate(deploymentMap["targetDeployDate"]) {
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
	if client == nil {
		a.clientGetterLock.Lock()
		client = a.mongoClient
		if client == nil {
			client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
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

func queryDeviceFirmwareInformation(db *mongo.Database, devId string) (error, map[string]string, map[string]string) {
	deploymentsColl := db.Collection("firmwareDeployments")

	dpBsonFilter := buildBsonFilter(devId)
	deploymentsBytes, err := deploymentsColl.FindOne(context.Background(), dpBsonFilter).DecodeBytes()
	if err != nil {
		log.Printf("Connection query firmware error: %v", err)
	}
	if len(deploymentsBytes) == 0 {
		return nil, nil, nil
	}

	deploymentMap := make(map[string]string)
	json.Unmarshal(deploymentsBytes, &deploymentMap)

	firmwareVersion := deploymentMap["firmwareVersion"]
	firmwareVersionColl := db.Collection("firmwareVersions")
	fvBsonFilter := buildFirmwareVersionFilter(firmwareVersion)
	firmwareVersionsBytes, err := firmwareVersionColl.FindOne(context.Background(), fvBsonFilter).DecodeBytes()
	if err != nil {
		log.Printf("Connection query firmware versions error: %v", err)
	}

	firmwareVersionMap := make(map[string]string)
	json.Unmarshal(firmwareVersionsBytes, &firmwareVersionMap)

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

func handleUpgradeCommand(firmwareVersionMap map[string]string) []byte {
	// firmwareName := []byte(firmwareVersionMap["firmwareName"])
	firmwareLength, _ := strconv.Atoi(firmwareVersionMap["fileSize"])
	crcValue, _ := strconv.Atoi(firmwareVersionMap["crcValue"])
	firmwareVersion := []byte(firmwareVersionMap["firmwareVersion"])
	operator := []byte(firmwareVersionMap["createdBy"])
	firmwareDesc := []byte(firmwareVersionMap["firmwareDesc"])
	identifier := []byte(firmwareVersionMap["identifier"])

	serialNumber := []byte{0x00}

	contentLength := make([]byte, 2)
	binary.BigEndian.PutUint16(contentLength, 118)

	// download firmware length
	downloadFirmwareLength := make([]byte, 4)
	binary.BigEndian.PutUint16(downloadFirmwareLength, uint16(firmwareLength))

	// crcValue
	crcValueTag := make([]byte, 2)
	binary.BigEndian.PutUint16(crcValueTag, uint16(crcValue))

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
	_, err := firmwareDeploymentColl.UpdateOne(context.Background(), bson.M{"devId": deviceId}, valueMap)
	if err != nil {
		return err
	}
	return nil
}
