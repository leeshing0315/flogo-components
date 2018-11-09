package mongodbpool

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_URI  = "mongodb://lima-w10:27017"
	TEST_DB   = "test"
	TEST_COLL = "items"
)

var coll *mongo.Collection

func init() {
	//todo implement shared sessions
	// client, err := mongo.NewClient(TEST_URI)

	//To remove below  error:
	// data not inserted topology is closed

	client, err := mongo.Connect(context.Background(), TEST_URI, nil)
	if err != nil {
		// warn and skip tests
	}

	db := client.Database(TEST_DB)
	coll = db.Collection(TEST_COLL)
}

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

func randomString(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func insert(dataVal interface{}) (interface{}, error) {

	result, err := coll.InsertOne(
		context.Background(),
		dataVal,
	)
	if err != nil {
		return nil, err
	}

	logger.Debug("Inserted: ", result.InsertedID)

	return result.InsertedID, nil
}

func delete(id interface{}) {
	oid := id.(objectid.ObjectID)
	_, err := coll.DeleteOne(context.Background(), bson.NewDocument(bson.EC.ObjectID("_id", oid)))
	if err != nil {
		logger.Debugf("Error Deleting [%s] : %s", id, err.Error())
		return
	}
	logger.Debug("Deleted", id)
}

// TestDelete
func TestDelete(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	name := randomString(5)
	val := map[string]interface{}{"name": name, "value": "blah"}
	_, err := insert(val)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	tc.SetInput("uri", TEST_URI)
	tc.SetInput("dbName", TEST_DB)
	tc.SetInput("collection", TEST_COLL)
	tc.SetInput("method", `DELETE`)

	tc.SetInput(ivKeyName, "name")
	tc.SetInput(ivKeyValue, name)

	_, deleteErr := act.Eval(tc)
	if deleteErr != nil {
		t.Error("data not deleted")
		t.Fail()
	}
}

// TestInsert
func TestInsert(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	tc.SetInput("uri", TEST_URI)
	tc.SetInput("dbName", TEST_DB)
	tc.SetInput("collection", TEST_COLL)
	tc.SetInput("method", `INSERT`)

	name := randomString(5)
	val := map[string]interface{}{"name": name, "value1": "foo", "value2": "foo2"}
	tc.SetInput(ivData, val)

	_, insertErr := act.Eval(tc)
	if insertErr != nil {
		t.Error("data not inserted", insertErr)
		t.Fail()
	}

	delete(tc.GetOutput(ovOutput))
}

func TestConcurrentInsert(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	ch := make(chan interface{})

	time1 := time.Now().Unix()
	for i := 0; i < 20000; i++ {
		go concurrentInsert(t, &act, ch)
	}

	for i := 0; i < 20000; i++ {
		<-ch
		// delete(<-ch)
	}

	time2 := time.Now().Unix()
	fmt.Println(time2 - time1)
}

func concurrentInsert(t *testing.T, act *activity.Activity, ch chan interface{}) {
	tc := test.NewTestActivityContext(getActivityMetadata())
	tc.SetInput("uri", TEST_URI)
	tc.SetInput("dbName", TEST_DB)
	tc.SetInput("collection", TEST_COLL)
	tc.SetInput("method", `INSERT`)

	name := randomString(5)
	val := map[string]interface{}{
		"battery":      81.0,
		"ReportType":   "例行数据",
		"outerTemp":    32.0,
		"innterTemp":   3276.7,
		"imeiNo":       "030146431",
		"dataTime":     "2018-09-28 07:25:18",
		"latitude":     22.558838,
		"longitude":    114.243935,
		"equipmentNo":  "CSNU7219310",
		"deviceNo":     "a770d84dff92423090db79c4fb976943",
		"location":     "广东省深圳市盐田区海景路" + name,
		"receivedTime": "2018-11-08 04:03:43.903435",
	}
	tc.SetInput(ivData, val)

	_, insertErr := (*act).Eval(tc)
	if insertErr != nil {
		t.Error("data not inserted", insertErr)
		t.Fail()
	}

	ch <- tc.GetOutput(ovOutput)
}

// TestReplace
func TestReplace(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	name := randomString(5)
	val := map[string]interface{}{"name": name, "value1": "foo", "value2": "foo2"}
	id, err := insert(val)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	tc.SetInput("uri", TEST_URI)
	tc.SetInput("dbName", TEST_DB)
	tc.SetInput("collection", TEST_COLL)
	tc.SetInput("method", `REPLACE`)

	val2 := map[string]interface{}{"name": name, "value1": "bar1", "value2": "bar2"}

	tc.SetInput(ivKeyName, "name")
	tc.SetInput(ivKeyValue, name)
	tc.SetInput(ivData, val2)

	_, replaceErr := act.Eval(tc)
	if replaceErr != nil {
		t.Error("data not replaced", replaceErr)
		t.Fail()
	}

	delete(id)
}

// TestReplace
func TestUpdate(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	name := randomString(5)
	val := map[string]interface{}{"name": name, "value1": "foo", "value2": "foo2"}
	id, err := insert(val)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	tc.SetInput("uri", TEST_URI)
	tc.SetInput("dbName", TEST_DB)
	tc.SetInput("collection", TEST_COLL)
	tc.SetInput("method", `UPDATE`)

	val2 := map[string]interface{}{"name": name, "value1": "bar1"}

	tc.SetInput(ivKeyName, "name")
	tc.SetInput(ivKeyValue, name)
	tc.SetInput(ivData, val2)

	_, updateErr := act.Eval(tc)
	if updateErr != nil {
		t.Error("update error", updateErr)
		t.Fail()
	}

	delete(id)
}

// TestGet
func TestGet(t *testing.T) {
	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	name := randomString(5)
	val := map[string]interface{}{"name": name, "value1": "foo", "value2": "foo2"}
	id, err := insert(val)
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	tc.SetInput("uri", TEST_URI)
	tc.SetInput("dbName", TEST_DB)
	tc.SetInput("collection", TEST_COLL)
	tc.SetInput("method", `GET`)

	tc.SetInput(ivKeyName, "name")
	tc.SetInput(ivKeyValue, name)

	_, getErr := act.Eval(tc)
	if getErr != nil {
		t.Error("unable to get data", getErr)
		t.Fail()
	}

	result := tc.GetOutput(ovOutput)
	assert.NotNil(t, result)

	delete(id)
}
