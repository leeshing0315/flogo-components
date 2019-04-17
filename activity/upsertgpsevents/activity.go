package upsertgpsevents

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-upsertgpsevents")

const (
	ivConnectionURI = "uri"
	ivDbName        = "dbName"

	ivGpsEvents     = "gpsevents"
	ivOpmodeChanges = "operationmodechanges"
	ivDeviceErrors  = "deviceerrors"

	ovOutput = "output"
	ovCount  = "count"
)

func init() {
	activityLog.SetLogLevel(logger.InfoLevel)
}

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

	// do eval
	connectionURI, _ := ctx.GetInput(ivConnectionURI).(string)
	dbName, _ := ctx.GetInput(ivDbName).(string)
	gpsEventsVal := ctx.GetInput(ivGpsEvents)
	opmodeChangesVal := ctx.GetInput(ivOpmodeChanges)
	deviceErrorsVal := ctx.GetInput(ivDeviceErrors)

	client := a.mongoClient
	if client == nil {
		a.clientGetterLock.Lock()
		client = a.mongoClient
		if client == nil {
			client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(connectionURI))
			if err != nil {
				activityLog.Errorf("Connection error: %v", err)
				a.clientGetterLock.Unlock()
				return false, err
			}
			a.mongoClient = client
		}
		a.clientGetterLock.Unlock()
	}

	db := client.Database(dbName)

	gpsEventsColl := db.Collection("gpsevents")
	if gpsEventsVal == nil || gpsEventsVal.(string) == "" {
		return true, nil
	}
	if strings.HasPrefix(gpsEventsVal.(string), "[") {
		var valueArray []string
		err = json.Unmarshal([]byte(gpsEventsVal.(string)), &valueArray)
		if err != nil {
			return false, err
		}
		var insertedIDArray []string
		var resultArray = []map[string]interface{}{}
		for _, val := range valueArray {
			var valueMap = make(map[string]interface{})
			err = json.Unmarshal([]byte(val), &valueMap)
			if err != nil {
				return false, err
			}
			result, err := gpsEventsColl.UpdateOne(
				context.Background(),
				map[string]interface{}{
					"cntrNum": valueMap["cntrNum"].(string),
					"seqno":   valueMap["seqno"].(string),
					"cltTime": valueMap["cltTime"].(string),
				},
				valueMap,
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return false, err
			}
			activityLog.Debugf("Insert Results $#v", result)
			insertedIDArray = append(insertedIDArray, result.UpsertedID.(primitive.ObjectID).String())
			valueMap["id"] = result.UpsertedID.(primitive.ObjectID).Hex()
			resultArray = append(resultArray, valueMap)
		}
		ctx.SetOutput(ovOutput, strings.Join(insertedIDArray, ","))
		resultArrayBytes, err := json.Marshal(resultArray)
		if err == nil {
			ctx.SetOutput("resultArray", string(resultArrayBytes))
		}
	} else {
		var valueMap = make(map[string]interface{})
		err = json.Unmarshal([]byte(gpsEventsVal.(string)), &valueMap)
		if err != nil {
			return false, err
		}
		result, err := gpsEventsColl.UpdateOne(
			context.Background(),
			map[string]interface{}{
				"cntrNum": valueMap["cntrNum"].(string),
				"seqno":   valueMap["seqno"].(string),
				"cltTime": valueMap["cltTime"].(string),
			},
			valueMap,
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return false, err
		}
		activityLog.Debugf("Insert Results $#v", result)
		ctx.SetOutput(ovOutput, result.UpsertedID)
		var resultArray = []map[string]interface{}{}
		valueMap["id"] = result.UpsertedID.(primitive.ObjectID).Hex()
		resultArray = append(resultArray, valueMap)
		resultArrayBytes, err := json.Marshal(resultArray)
		if err == nil {
			ctx.SetOutput("resultArray", string(resultArrayBytes))
		}
	}

	opmodeChangesColl := db.Collection("operationmodechanges")
	if opmodeChangesVal == nil || opmodeChangesVal.(string) == "" {
		return true, nil
	}
	if strings.HasPrefix(opmodeChangesVal.(string), "[") {
		var valueArray []string
		err = json.Unmarshal([]byte(opmodeChangesVal.(string)), &valueArray)
		if err != nil {
			return false, err
		}
		for _, val := range valueArray {
			var valueMap = make(map[string]interface{})
			err = json.Unmarshal([]byte(val), &valueMap)
			if err != nil {
				return false, err
			}
			result, err := opmodeChangesColl.UpdateOne(
				context.Background(),
				map[string]interface{}{
					"cntrnum": valueMap["cntrnum"].(string),
					"seqno":   valueMap["seqno"].(string),
					"logtime": valueMap["logtime"].(string),
					"opmode":  valueMap["opmode"].(string),
				},
				valueMap,
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return false, err
			}
			activityLog.Debugf("Insert Results $#v", result)
		}
	} else {
		var valueMap = make(map[string]interface{})
		err = json.Unmarshal([]byte(opmodeChangesVal.(string)), &valueMap)
		if err != nil {
			return false, err
		}
		result, err := opmodeChangesColl.UpdateOne(
			context.Background(),
			map[string]interface{}{
				"cntrnum": valueMap["cntrnum"].(string),
				"seqno":   valueMap["seqno"].(string),
				"logtime": valueMap["logtime"].(string),
				"opmode":  valueMap["opmode"].(string),
			},
			valueMap,
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return false, err
		}
		activityLog.Debugf("Insert Results $#v", result)
	}

	deviceErrorsColl := db.Collection("deviceerrors")
	if deviceErrorsVal == nil || deviceErrorsVal.(string) == "" {
		return true, nil
	}
	if strings.HasPrefix(deviceErrorsVal.(string), "[") {
		var valueArray []string
		err = json.Unmarshal([]byte(deviceErrorsVal.(string)), &valueArray)
		if err != nil {
			return false, err
		}
		for _, val := range valueArray {
			var valueMap = make(map[string]interface{})
			err = json.Unmarshal([]byte(val), &valueMap)
			if err != nil {
				return false, err
			}
			result, err := deviceErrorsColl.UpdateOne(
				context.Background(),
				map[string]interface{}{
					"cntrnum": valueMap["cntrnum"].(string),
					"seqno":   valueMap["seqno"].(string),
					"logtime": valueMap["logtime"].(string),
					"opmode":  valueMap["opmode"].(string),
				},
				valueMap,
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return false, err
			}
			activityLog.Debugf("Insert Results $#v", result)
		}
	} else {
		var valueMap = make(map[string]interface{})
		err = json.Unmarshal([]byte(deviceErrorsVal.(string)), &valueMap)
		if err != nil {
			return false, err
		}
		result, err := deviceErrorsColl.UpdateOne(
			context.Background(),
			map[string]interface{}{
				"cntrnum": valueMap["cntrnum"].(string),
				"seqno":   valueMap["seqno"].(string),
				"logtime": valueMap["logtime"].(string),
				"opmode":  valueMap["opmode"].(string),
			},
			valueMap,
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return false, err
		}
		activityLog.Debugf("Insert Results $#v", result)
	}

	return true, nil
}
