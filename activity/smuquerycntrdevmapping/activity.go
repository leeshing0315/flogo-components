package smuquerycntrdevmapping

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/leeshing0315/flogo-components/common/entity"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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

	// do eval
	pin := ctx.GetInput("pin").(string)
	eventTime := ctx.GetInput("eventTime").(string)
	originalPacket := ctx.GetInput("originalPacket").([]byte)

	client := a.mongoClient
	if client == nil {
		a.clientGetterLock.Lock()
		client = a.mongoClient
		if client == nil {
			client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
			if err != nil {
				log.Printf("Connection error: %v", err)
				a.clientGetterLock.Unlock()
				return false, err
			}
			a.mongoClient = client
		}
		a.clientGetterLock.Unlock()
	}

	db := client.Database("iot")

	coll := db.Collection("containerDeviceMappings")

	doc := make(map[string]interface{})
	doc["pin"] = pin
	doc["status"] = "active"

	result := coll.FindOne(context.Background(), doc)
	val := make(map[string]interface{})
	err = result.Decode(&val)

	// if no pin-mapping, log down to DB
	if err != nil {
		exceptionColl := db.Collection("tcpServerExceptions")
		exception := &entity.TcpServerException{
			RevTime:      eventTime,
			Pin:          pin,
			ProtocolType: entity.PROTOCOL_TYPE_SMU,
			Bytes:        originalPacket,
			ErrorType:    entity.ERROR_TYPE_PIN_NOT_REGISTERED,
			ErrorReason:  "Pin has not registered",
		}
		exceptionColl.InsertOne(context.Background(), exception)
		return false, err
	}

	cntrDevMappingBytes, err := json.Marshal(val)
	if err != nil {
		return false, err
	}

	ctx.SetOutput("output", string(cntrDevMappingBytes))

	return true, nil
}
