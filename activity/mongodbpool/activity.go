package mongodbpool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-flogo-mongodb-pool")

const (
	methodGet     = "GET"
	methodDelete  = "DELETE"
	methodInsert  = "INSERT"
	methodReplace = "REPLACE"
	methodUpdate  = "UPDATE"

	ivConnectionURI = "uri"
	ivDbName        = "dbName"
	ivCollection    = "collection"
	ivMethod        = "method"

	ivKeyName  = "keyName"
	ivKeyValue = "keyValue"
	ivData     = "data"

	ovOutput = "output"
	ovCount  = "count"
)

func init() {
	activityLog.SetLogLevel(logger.InfoLevel)
}

/*
Integration with MongoDb
inputs: {uri, dbName, collection, method, [keyName, keyValue, value]}
outputs: {output, count}
*/
type MongoDbActivity struct {
	metadata         *activity.Metadata
	mongoClient      *mongo.Client
	clientGetterLock sync.Mutex
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MongoDbActivity{metadata: metadata}
}

// Metadata returns the activity's metadata
func (a *MongoDbActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements api.Activity.Eval - MongoDb integration
func (a *MongoDbActivity) Eval(ctx activity.Context) (done bool, err error) {

	//mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
	connectionURI, _ := ctx.GetInput(ivConnectionURI).(string)
	dbName, _ := ctx.GetInput(ivDbName).(string)
	collectionName, _ := ctx.GetInput(ivCollection).(string)
	method, _ := ctx.GetInput(ivMethod).(string)
	keyName, _ := ctx.GetInput(ivKeyName).(string)
	keyValue, _ := ctx.GetInput(ivKeyValue).(string)
	value := ctx.GetInput(ivData)

	client := a.mongoClient
	if client == nil {
		a.clientGetterLock.Lock()
		client = a.mongoClient
		if client == nil {
			client, err = mongo.Connect(context.Background(), connectionURI, nil)
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

	coll := db.Collection(collectionName)

	switch strings.ToUpper(method) {
	case methodGet:
		document, buildErr := buildDocument(keyName, keyValue)
		if buildErr != nil {
			return false, buildErr
		}
		result := coll.FindOne(context.Background(), document)
		val := make(map[string]interface{})
		err := result.Decode(val)
		if err != nil {
			return false, err
		}

		activityLog.Debugf("Get Results $#v", result)

		ctx.SetOutput(ovOutput, val)
	case methodDelete:
		document, buildErr := buildDocument(keyName, keyValue)
		if buildErr != nil {
			return false, buildErr
		}
		result, err := coll.DeleteMany(
			context.Background(),
			document,
		)
		if err != nil {
			return false, err
		}

		activityLog.Debugf("Delete Results $#v", result)

		ctx.SetOutput(ovCount, result.DeletedCount)
	case methodInsert:
		if value.(string) == "" {
			break
		}
		if strings.HasPrefix(value.(string), "[") {
			var valueArray []string
			err = json.Unmarshal([]byte(value.(string)), &valueArray)
			if err != nil {
				return false, err
			}
			var insertedIDArray []string
			for _, val := range valueArray {
				var valueMap map[string]interface{}
				err = json.Unmarshal([]byte(val), &valueMap)
				if err != nil {
					return false, err
				}
				result, err := coll.InsertOne(
					context.Background(),
					valueMap,
				)
				if err != nil {
					return false, err
				}
				activityLog.Debugf("Insert Results $#v", result)
				insertedIDArray = append(insertedIDArray, result.InsertedID.(primitive.ObjectID).String())
			}
			ctx.SetOutput(ovOutput, strings.Join(insertedIDArray, ","))
		} else {
			var valueMap map[string]interface{}
			err = json.Unmarshal([]byte(value.(string)), &valueMap)
			if err != nil {
				return false, err
			}
			result, err := coll.InsertOne(
				context.Background(),
				valueMap,
			)
			if err != nil {
				return false, err
			}
			activityLog.Debugf("Insert Results $#v", result)
			ctx.SetOutput(ovOutput, result.InsertedID)
		}
	case methodReplace:
		document, buildErr := buildDocument(keyName, keyValue)
		if buildErr != nil {
			return false, buildErr
		}
		var valueMap map[string]interface{}
		err = json.Unmarshal([]byte(value.(string)), &valueMap)
		if err != nil {
			return false, err
		}
		result, err := coll.ReplaceOne(
			context.Background(),
			document,
			valueMap,
		)
		if err != nil {
			return false, err
		}

		activityLog.Debugf("Replace Results $#v", result)
		ctx.SetOutput(ovOutput, result.UpsertedID)
		ctx.SetOutput(ovCount, result.ModifiedCount)

	case methodUpdate:
		document, buildErr := buildDocument(keyName, keyValue)
		if buildErr != nil {
			return false, buildErr
		}
		var valueMap map[string]interface{}
		err = json.Unmarshal([]byte(value.(string)), &valueMap)
		if err != nil {
			return false, err
		}
		result, err := coll.UpdateOne(
			context.Background(),
			document,
			bson.M{"$set": valueMap},
		)
		if err != nil {
			return false, err
		}

		activityLog.Debugf("Update Results $#v", result)
		ctx.SetOutput(ovOutput, result.UpsertedID)
		ctx.SetOutput(ovCount, result.ModifiedCount)
	default:
		activityLog.Errorf("unsupported method '%s'", method)
		return false, fmt.Errorf("unsupported method '%s'", method)
	}

	return true, nil
}

func buildDocument(keyName string, keyValue string) (interface{}, error) {
	names := strings.Split(keyName, ",")
	values := strings.Split(keyValue, ",")

	namesLen := len(names)
	valuesLen := len(values)
	if namesLen != valuesLen {
		return nil, errors.New("KeyValueLenNotMatch")
	}

	result := bson.M{}

	for i := 0; i < namesLen; i++ {
		result[names[i]] = values[i]
	}

	return result, nil
}
