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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ActivityLog is the default logger for the Log Activity
var activityLog = logger.GetLogger("activity-flogo-mongodb-pool")

const (
	methodGet        = "GET"
	methodGetMany    = "GETMANY"
	methodDelete     = "DELETE"
	methodInsert     = "INSERT"
	methodReplace    = "REPLACE"
	methodUpdate     = "UPDATE"
	methodUpdateMany = "UPDATEMANY"

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

	coll := db.Collection(collectionName)

	switch strings.ToUpper(method) {
	case methodGet:
		document, buildErr := buildDocument(keyName, keyValue)
		if buildErr != nil {
			return false, buildErr
		}
		result := coll.FindOne(context.Background(), document)
		val := make(map[string]interface{})
		err := result.Decode(&val)
		if err != nil {
			return false, err
		}

		activityLog.Debugf("Get Results $#v", result)

		ctx.SetOutput(ovOutput, val)
	case methodGetMany:
		condition, buildErr := buildDocument(keyName, keyValue)
		if buildErr != nil {
			return false, buildErr
		}
		cur, queryErr := coll.Find(context.Background(), condition)
		if queryErr != nil {
			return false, queryErr
		}
		defer cur.Close(context.Background())
		var count = 0
		result := []map[string]interface{}{}
		for cur.Next(context.Background()) {
			elem := make(map[string]interface{})
			if decodeErr := cur.Decode(&elem); decodeErr != nil {
				return false, decodeErr
			}
			result = append(result, elem)
			count++
		}
		if err := cur.Err(); err != nil {
			return false, err
		}
		ctx.SetOutput(ovCount, count)
		if count == 0 {
			activityLog.Debugf("No Document can be found")
			ctx.SetOutput(ovOutput, make(map[string]interface{}))
		} else {
			activityLog.Debugf("Get Multiple Results $#v", result)
			ctx.SetOutput(ovOutput, result)
		}
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
		if value == nil || value.(string) == "" {
			break
		}
		if strings.HasPrefix(value.(string), "[") {
			var valueArray []string
			err = json.Unmarshal([]byte(value.(string)), &valueArray)
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
				result, err := coll.InsertOne(
					context.Background(),
					valueMap,
				)
				if err != nil {
					return false, err
				}
				activityLog.Debugf("Insert Results $#v", result)
				insertedIDArray = append(insertedIDArray, result.InsertedID.(primitive.ObjectID).String())
				valueMap["id"] = result.InsertedID.(primitive.ObjectID).Hex()
				resultArray = append(resultArray, valueMap)
			}
			ctx.SetOutput(ovOutput, strings.Join(insertedIDArray, ","))
			resultArrayBytes, err := json.Marshal(resultArray)
			if err == nil {
				ctx.SetOutput("resultArray", string(resultArrayBytes))
			}
		} else {
			var valueMap = make(map[string]interface{})
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
			var resultArray = []map[string]interface{}{}
			valueMap["id"] = result.InsertedID.(primitive.ObjectID).Hex()
			resultArray = append(resultArray, valueMap)
			resultArrayBytes, err := json.Marshal(resultArray)
			if err == nil {
				ctx.SetOutput("resultArray", string(resultArrayBytes))
			}
		}
	case methodReplace:
		if value == nil || value.(string) == "" {
			break
		}
		document, buildErr := buildDocument(keyName, keyValue)
		if buildErr != nil {
			return false, buildErr
		}
		var valueMap = make(map[string]interface{})
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
		if value == nil || value.(string) == "" {
			break
		}
		document, buildErr := buildDocument(keyName, keyValue)
		if buildErr != nil {
			return false, buildErr
		}
		var valueMap = make(map[string]interface{})
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
	case methodUpdateMany:
		if value == nil || value.(string) == "" {
			break
		}
		document, buildErr := buildDocument(keyName, keyValue)
		if buildErr != nil {
			return false, buildErr
		}
		var valueMap = make(map[string]interface{})
		err = json.Unmarshal([]byte(value.(string)), &valueMap)
		if err != nil {
			return false, err
		}
		result, err := coll.UpdateMany(
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
	if keyName != "" {
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
	} else {
		result := make(map[string]interface{})
		err := json.Unmarshal([]byte(keyValue), &result)
		if err != nil {
			return nil, errors.New("KeyValueNotJson")
		}

		return result, nil
	}
}
