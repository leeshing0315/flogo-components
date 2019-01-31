package smuautoreg

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
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
	autoReg := ctx.GetInput("autoReg").(string)
	pin := ctx.GetInput("pin").(string)
	cntrNum := ctx.GetInput("cntrNum").(string)
	devId := ctx.GetInput("devId").(string)

	// if not autoReg, just pass
	if autoReg != "true" {
		return true, nil
	}

	// do eval
	client, err := mongo.Connect(context.Background(), uri, nil)
	defer client.Disconnect(context.Background())
	if err != nil {
		log.Printf("Connection error: %v", err)
		return true, nil
	}

	db := client.Database(dbName)

	coll := db.Collection("containerDeviceMappings")

	// find old active
	filter := buildFilter(pin)
	oldActiveBytes, err := coll.FindOne(context.Background(), filter).DecodeBytes()
	if err != nil {
		log.Printf("FindOne error: %v", err)
		return true, nil
	}
	if len(oldActiveBytes) == 0 {
		return
	}
	oldActive := make(map[string]interface{})
	err = json.Unmarshal(oldActiveBytes, &oldActive)
	if err != nil {
		log.Printf("FindOne error: %v", err)
		return true, nil
	}

	// if not change, just pass
	if isNumberNotChanged(oldActive, cntrNum, devId) {
		log.Printf("AutoReg not changed")
		return true, nil
	}

	// update old from active to inactive
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	cntrDevMappingDateStr := time.Now().In(loc).Format("2006-01-02 15:04:05")
	update := buildUpdate(cntrDevMappingDateStr)
	_, err = coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("UpdateOne error: %v", err)
		return true, nil
	}

	// insert new active
	newActive := buildNewActive(oldActive, cntrNum, devId, cntrDevMappingDateStr)
	insertResult, err := coll.InsertOne(context.Background(), newActive)
	if err != nil {
		log.Printf("InsertOne error: %v", err)
		coll.UpdateOne(
			context.Background(),
			map[string]interface{}{"_id": oldActive["_id"]},
			map[string]interface{}{"status": "active"},
		)
		return true, nil
	}
	newActive["_id"] = insertResult.InsertedID.(primitive.ObjectID).Hex()

	// insert auditLog
	auditLog := buildAuditLog(oldActive, newActive)
	deviceAuditLogsColl := db.Collection("deviceAuditLogs")
	deviceAuditLogsColl.InsertOne(context.Background(), auditLog)

	return true, nil
}

func buildFilter(pin string) map[string]interface{} {
	filter := make(map[string]interface{})
	filter["pin"] = pin
	filter["status"] = "active"
	return filter
}

func buildUpdate(cntrDevMappingDateStr string) map[string]interface{} {
	update := make(map[string]interface{})
	update["status"] = "inactive"
	update["changetime"] = cntrDevMappingDateStr
	return update
}

func buildNewActive(oldActive map[string]interface{}, cntrNum, devId, cntrDevMappingDateStr string) map[string]interface{} {
	newActive := make(map[string]interface{})
	for k, v := range oldActive {
		newActive[k] = v
	}
	delete(newActive, "_id")
	delete(newActive, "id")
	newActive["status"] = "active"
	newActive["carno"] = "devId"
	newActive["carid"] = "cntrNum"
	newActive["regtime"] = cntrDevMappingDateStr
	newActive["changetime"] = cntrDevMappingDateStr
	newActive["lastUpdated"] = cntrDevMappingDateStr
	return newActive
}

func buildAuditLog(oldActive, newActive map[string]interface{}) map[string]interface{} {
	auditLog := make(map[string]interface{})
	auditLog["action"] = "inactive"
	auditLog["createdBy"] = "TCP_SERVER"
	auditLog["beforeValue"] = oldActive
	auditLog["afterValue"] = newActive
	auditLog["createdDate"] = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	return auditLog
}

func isNumberNotChanged(oldActive map[string]interface{}, cntrNum, devId string) bool {
	if oldActive["carid"] != cntrNum || oldActive["carno"] != devId {
		return true
	} else {
		return false
	}
}
