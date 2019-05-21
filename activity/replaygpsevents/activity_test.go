package replaygpsevents

import (
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
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
	decoded, err := base64.StdEncoding.DecodeString(`NwAHAIg2AAMAGwAAIDKZCSgZEjUDMMlVAJe8MAAFAABkBAIODNWyNgAEABsAACAymQkoGScxAzDFQQCXuh0ACQAAZAQCDgyF5zYABQAbAAAgMpkJKBlCMQMwzZsAl7lEAAUAAGMEAg4LLRI2AAYAGwAAADIZBRQZV0gDMMErAJfBgwAbAUNiBAIQCzovulQ=`)
	if err != nil {
		return
	}
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

	service.GenGpsEventFromSinglePacket(singlePacket, "seqNo", "cntrNum", "revTime", "COSU")
}
