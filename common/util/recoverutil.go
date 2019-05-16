package util

import (
	"context"
	"fmt"
	"os"

	"github.com/leeshing0315/flogo-components/common/entity"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var coll *mongo.Collection

func init() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		panic(err)
	}
	mongoClient = client
	mongoDb := client.Database("iot")
	coll = mongoDb.Collection("tcpServerExceptions")
}

func LogDownException(eventTime string, pin string, originalPacket []byte, err interface{}) {
	exception := &entity.TcpServerException{
		RevTime:      eventTime,
		Pin:          pin,
		ProtocolType: entity.PROTOCOL_TYPE_SMU,
		Bytes:        originalPacket,
		ErrorType:    entity.ERROR_TYPE_UNCHECKED,
		ErrorReason:  fmt.Sprint(err),
	}
	coll.InsertOne(context.Background(), exception)
}
