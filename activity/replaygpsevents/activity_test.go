package replaygpsevents

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/leeshing0315/flogo-components/common/entity"
	"github.com/leeshing0315/flogo-components/common/service"
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
	tc.SetInput("uri", "mongodb://lima-w10:27017")
	tc.SetInput("dbName", "test")
	tc.SetInput("reqDataSegment", []byte(`{"from":"2019-03-01","to":"2019-05-01"}`))

	act.Eval(tc)

	//check result attr
}

func TestSth(t *testing.T) {
	decoded, err := base64.StdEncoding.DecodeString(`NwlGA2o2CToAHgCAgHIZBgIVQAAAx0HoBgOiRwABADhkBAIPFAUBwPNxNgk7AEkAAAByGQYCFUlJAMdB6wYDomMAAAA4ZAQCEBQCLEQB////wPnD+ev9K/9X////JhAO8gNNAx3JSwH9wPcCygA/AUgAAP8N/5ZLgSU2CTwASQAAAHIZBgIVWUkAx0I4BgOiowAAADhkBAIPFAIsRAH////A+cP6Q8G3/1v///8mEw75A14DJ8kPAWbEIwMbAMkBSAAA/xsASEt3NTYJPQAeAICAchkGAhYEWQDHQcoGA6KPAAAAOGQEAg8UBQHBBYE2CT4ASQAAAHIZBgIWCUkAx0HaBgOidQAAADhkBAIQFAIsRAH////B+cP79/0r/4D///8mCQ7sBK4DwckX/zv7wwN3/20BSABH/1//hEuYzzYJPwBJAAAAchkGAhYZSQDHQgoGA6J3AAAAOGQEAg8UAixEAf///8H5w/qj++v/b////yYEDt8EDAMXyRf+5foLA4r/jgFIADf/K/9ZS+EoNglAAEkAAAByGQYCFilJAMdB1wYDonUAAAA4ZAQCEBQCLEQB////wfnD+jv7j/9w////JgQO0QQdAyfIz/7l+XsDSP9kAUgANv8a/0xLhCI2CUEASQAAAHIZBgIWOUkAx0HuBgOibQAAADhkBAIQFAIsRAH////B+cP59/tX/2v///8mBQ7fBAEDDciP/uf5IwMy/zEBSAA2/w//REt4SzYJQgBJAAAAchkGAhZJSQDHQZoGA6KFAAAAOGQEAhAUAixEAf///8H5w/n3+zf/bv///yYFDt8EOQNMyI/+2fkrAyr/SAFIADX/D/8+SzCHNglDAEkAAAByGQYCFllJAMdBwAYDonMAAAA4ZAQCEBMCLEQB////wfnD+eP7F7Ru////JgQO5QQMAxfI2/7K+UMDOv9/AUgANf8M/zpLmGQ2CUQASQAAAHIZBgIXCUkAx0GyBgOiewAAAStkBAIREwIsRAH////B+cP51/sDr2b///8mAw75A9kC4sjP/sv5GwNQ/2EBSAA0/wr/N0sBljYJRQBJAAAAchkGAhcZSQDHQd4GA6KYAAAAFWQEAhESAixEAf///8H5w/m7+vetbf///yYEDtgEAQMSyOP+2/jvAw//LgFIADT/Bv81S6OseMM=`)
	if err != nil {
		return
	}
	log.Println(convertBytesToStrings(decoded))
	_, reqDataSegment := parseBytes(decoded)
	packets := splitPackets(reqDataSegment)
	for _, dateSegment := range packets {
		singlePacket, err := service.ParseToSinglePacket(dateSegment)
		if err != nil {
			return
		}

		service.GenGpsEventFromSinglePacket(singlePacket, "seqNo", "cntrNum", "revTime", "COSU")
	}
}

func TestSinglePacketSth(t *testing.T) {
	decoded, err := base64.StdEncoding.DecodeString(`NgDJAFoAgAByGQUVCRFTAh/oRQc1jwoBDgCEZAQCEBICLEQB////w8M3wyPDV+w+////JyYQ1wTeAu3EiwB5wusCS//tACoAagB+AIP/CAJEXQkI/er96v3q/eoFAcP9kw==`)
	if err != nil {
		return
	}
	_, reqDataSegment := parseBytes(decoded)
	singlePacket, err := service.ParseToSinglePacket(reqDataSegment)
	if err != nil {
		return
	}
	// singlePacketBytes, _ := json.Marshal(singlePacket)
	// fmt.Println(string(singlePacketBytes))

	gpsEvent := service.GenGpsEventFromSinglePacket(singlePacket, "seqNo", "cntrNum", "revTime", "COSU")
	gpsEventBytes, _ := json.Marshal(gpsEvent)
	fmt.Println(string(gpsEventBytes))
}

func TestEventLog(t *testing.T) {
	decoded, err := base64.StdEncoding.DecodeString(`OALmATQBE5TlAOlGWBkJSlo8YU7//////wETlOTA6UZYGQlKWTxhTv//////AROTLAAyeM8ZAojPNilO//////8CE5Tk+6ACAhOU5PugAQITkywKIAEBE5TlQKgiDBuH8g9BWk7//////wITlOVAsacCE5TlQLGmAhOU5TKgBgITlOUyoAICE5TlDSACAhOU5QygBgITlOUMoAICE5TlDCACAhOU5QwA3AETlOWApY2TG4Zt60DPMv//////AhOU5XKwEAITlOVlsBECE5TlW7GjAhOU5VuxogITlOVbsaABE5TlwKMY+BuEJT1ALTL//////wETlOgAVA0cG4PZDzw6Tv//////AhOU5eogAgITlOXdsWwCE5Tl3bFrAhOU5d2wQQITlOXQsaUCE5Tl0LGkAhOU5dCxoXv9`)
	if err != nil {
		return
	}
	_, reqDataSegment := parseBytes(decoded)
	packets := splitEventLogPackets(reqDataSegment)
	var eventLogs []*entity.EventLog
	for _, dateSegment := range packets {
		eventLog := service.ParseToEventLog(dateSegment, time.Now(), "DUMMY", 0)
		eventLogs = append(eventLogs, eventLog)
	}
	bytes, _ := json.Marshal(eventLogs)
	println(string(bytes))
}

func splitEventLogPackets(data []byte) [][]byte {
	var result = [][]byte{}
	for i := data; len(i) > 1; {
		var dataSegmentLen int
		if i[0] == 1 {
			dataSegmentLen = 19
		} else if i[0] == 2 {
			dataSegmentLen = 6
		}
		dataSegment := i[0 : dataSegmentLen+1]
		result = append(result, dataSegment)

		i = i[dataSegmentLen+1:]
	}
	return result
}

func convertBytesToStrings(input []byte) []string {
	output := make([]string, len(input))
	for index, val := range input {
		output[index] = "0x" + strconv.FormatUint(uint64(val), 16)
	}
	return output
}
