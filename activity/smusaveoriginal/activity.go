package smusaveoriginal

import (
	"context"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/leeshing0315/flogo-components/common/entity"
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
	ip := ctx.GetInput("ip").(string)
	eventTime := ctx.GetInput("eventTime").(string)
	originalPacketBytes := ctx.GetInput("originalPacketBytes").([]byte)
	pin := ctx.GetInput("pin").(string)
	uri := ctx.GetInput("uri").(string)
	dbName := ctx.GetInput("dbName").(string)

	originalPacket := &entity.OriginalPacket{
		Protocol:    "SMU",
		Ip:          ip,
		RevTime:     eventTime,
		Bytes:       originalPacketBytes,
		BytesLength: len(originalPacketBytes),
		Pin:         pin,
		Source:      "TCP_SERVER",
	}
	if originalPacketBytes[0] == 0x32 {
		pinLen := originalPacketBytes[6]
		originalPacket.Pin = string(originalPacketBytes[7 : 7+pinLen])
	}

	client := a.mongoClient
	if client == nil {
		a.clientGetterLock.Lock()
		client = a.mongoClient
		if client == nil {
			client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
			if err != nil {
				a.clientGetterLock.Unlock()
				return false, err
			}
			a.mongoClient = client
		}
		a.clientGetterLock.Unlock()
	}

	db := client.Database(dbName)

	coll := db.Collection("originalPackets")

	_, err = coll.InsertOne(context.Background(), originalPacket)
	if err != nil {
		println(err)
		return false, err
	}

	return true, nil
}
