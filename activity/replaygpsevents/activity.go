package replaygpsevents

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"github.com/leeshing0315/flogo-components/common/entity"
	"github.com/leeshing0315/flogo-components/common/service"
)

var activityLog = logger.GetLogger("replay-gpsevents")

var loc *time.Location

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

var cntrDeviceMappings map[string]map[string]interface{} = make(map[string]map[string]interface{})

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(ctx activity.Context) (done bool, err error) {
	loc, _ = time.LoadLocation("Asia/Hong_Kong")

	// do eval
	uri, _ := ctx.GetInput("uri").(string)
	dbName, _ := ctx.GetInput("dbName").(string)
	reqDataSegment, _ := ctx.GetInput("reqDataSegment").([]byte)
	// eventTime, _ := ctx.GetInput("eventTime").(string)

	client := a.mongoClient
	if client == nil {
		a.clientGetterLock.Lock()
		client = a.mongoClient
		if client == nil {
			client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
			if err != nil {
				activityLog.Errorf("Connection error: %v", err)
				a.clientGetterLock.Unlock()
				return false, err
			}
			a.mongoClient = client
		}
		a.clientGetterLock.Unlock()
	}

	jsonMap := make(map[string]interface{})
	json.Unmarshal(reqDataSegment, &jsonMap)
	from := jsonMap["from"].(string)
	to := jsonMap["to"].(string)

	db := client.Database(dbName)

	err = loadAllCntrDeviceMappings(db)
	if err != nil {
		activityLog.Errorf("%v", err)
		return true, nil
	}

	err = handleOriginalPackets(db, from, to)
	if err != nil {
		activityLog.Errorf("%v", err)
		return true, nil
	}

	return true, nil
}

func loadAllCntrDeviceMappings(db *mongo.Database) error {
	coll := db.Collection("containerDeviceMappings")

	cursor, err := coll.Find(
		context.Background(),
		bson.M{},
	)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		mapping := make(map[string]interface{})
		err := cursor.Decode(&mapping)
		if err != nil {
			return err
		}

		cntrDeviceMappings[mapping["pin"].(string)] = mapping
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

func handleOriginalPackets(db *mongo.Database, from string, to string) error {
	coll := db.Collection("originalPackets")

	cursor, err := coll.Find(
		context.Background(),
		bson.M{
			"revtime": bson.M{
				"$gte": from,
				"$lt":  to,
			},
			"$or": bson.A{
				bson.M{"replayVersion": bson.M{"$exists": false}},
				bson.M{"replayVersion": bson.M{"$lt": 1}},
			},
		},
	)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		originalPacket := make(map[string]interface{})
		err := cursor.Decode(&originalPacket)
		if err != nil {
			return err
		}

		err = handleOneOriginalPacket(db, originalPacket)
		if err != nil {
			return err
		}
	}
	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

func parseBytes(bytes []byte) (seqNo string, reqDataSegment []byte) {
	seqNo = strconv.FormatUint(uint64(binary.BigEndian.Uint16(bytes[1:3])), 10)
	dataSegmentLen := int(binary.BigEndian.Uint16(bytes[3:5]))
	reqDataSegment = bytes[5 : 5+dataSegmentLen]
	return seqNo, reqDataSegment
}

func handleOneOriginalPacket(db *mongo.Database, originalPacket map[string]interface{}) (err error) {
	bytes := originalPacket["bytes"].([]byte)
	if bytes[0] != 0x36 && bytes[0] != 0x37 {
		return nil
	}
	// revTime := originalPacket["revtime"].(string)
	revTime := time.Now().In(loc).Format("2006-01-02T15:04:05+08:00")
	pin := originalPacket["pin"].(string)
	cntrNum := cntrDeviceMappings[pin]["carid"].(string)

	seqNo, reqDataSegment := parseBytes(bytes)

	switch bytes[0] {
	case 0x36:
		err = handleSinglePacket(db, seqNo, reqDataSegment, revTime, cntrNum)
		if err != nil {
			return err
		}
	case 0x37:
		err = handleMultiPacket(db, seqNo, reqDataSegment, revTime, cntrNum)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleSinglePacket(db *mongo.Database, seqNo string, reqDataSegment []byte, revTime string, cntrNum string) error {
	singlePacket, err := service.ParseToSinglePacket(reqDataSegment)
	if err != nil {
		return err
	}

	gpsEvent := service.GenGpsEventFromSinglePacket(singlePacket, seqNo, cntrNum, revTime)
	err = handleOneGpsEvent(db, gpsEvent)
	if err != nil {
		return err
	}

	return nil
}

func splitPackets(data []byte) [][]byte {
	var result = [][]byte{}
	for i := data; len(i) > 1; {
		dataSegmentLen := binary.BigEndian.Uint16(i[3:5])
		dataSegment := i[5 : dataSegmentLen+5]
		result = append(result, dataSegment)

		i = i[dataSegmentLen+7:]
	}
	return result
}
func handleMultiPacket(db *mongo.Database, seqNo string, reqDataSegment []byte, revTime string, cntrNum string) (err error) {
	packets := splitPackets(reqDataSegment)
	for _, dateSegment := range packets {
		err = handleSinglePacket(db, seqNo, dateSegment, revTime, cntrNum)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleOneGpsEvent(db *mongo.Database, gpsEvent *entity.GpsEvent) (err error) {
	gpsEventsColl := db.Collection("gpsevents")
	gpsEventMappingColl := db.Collection("gpseventMappings")
	// update gpsevent
	_, err = gpsEventsColl.UpdateOne(
		context.Background(),
		bson.M{
			"cntrNum": gpsEvent.CntrNum,
			"cltTime": gpsEvent.CltTime,
		},
		bson.M{
			"$set": bson.M{
				"eleState": gpsEvent.EleState,
				"posFlag":  gpsEvent.PosFlag,
				"lat":      gpsEvent.Lat,
				"lng":      gpsEvent.Lng,
			},
		},
	)
	if err != nil {
		return err
	}
	// update gpseventmapping
	_, err = gpsEventMappingColl.UpdateOne(
		context.Background(),
		bson.M{
			"cntrNum": gpsEvent.CntrNum,
			"cltTime": gpsEvent.CltTime,
		},
		bson.M{
			"$set": bson.M{
				"eleState": gpsEvent.EleState,
				"posFlag":  gpsEvent.PosFlag,
				"lat":      gpsEvent.Lat,
				"lng":      gpsEvent.Lng,
			},
		},
	)
	if err != nil {
		return err
	}
	return nil
}
