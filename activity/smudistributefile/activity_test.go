package smudistributefile

import (
	"context"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs
	tc.SetInput("firmwareVersion", `{"filePath": "./HO16.obj"}`)
	tc.SetInput("serialNumber", `2`)

	act.Eval(tc)

	//check result attr
}

func TestFileContent(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fail()
		return
	}
	db := client.Database("iot")
	coll := db.Collection("firmwareVersions")
	result := coll.FindOne(context.Background(), map[string]interface{}{"firmwareVersion": "HS19-2"})
	if err != nil {
		t.Fail()
		return
	}
	resultMap := make(map[string]interface{})
	result.Decode(&resultMap)
	bytes := getBytesFromMap(resultMap["fileContent"].(primitive.A)[0].(map[string]interface{}))
	log.Println(bytes)
	log.Println(time.Now().UTC().Format("2006-01-02T15:04:05Z"))
}
