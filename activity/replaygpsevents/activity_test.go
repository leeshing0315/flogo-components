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
	decoded, err := base64.StdEncoding.DecodeString(`NwAPAZg2AAMAGwAAADIZBQgSIwQB3ZPdB0CqhwAAAABkBAIOEcq3NgAEABsAAAAyGQUIEjgQAd2UQgdAqVsAAQAAZAQCDhE71jYABQAbAAAAMhkFCBJTEwHdk6oHQKmhAAAAAGQEAg0QeBo2AAYAGwAAADIZBQgTCBcB3ZOSB0CopwAAAABkBAIOEWwINgAHABsAAAAyGQUIEyMhAd2VWwdApgYAAwAAZAQCEBHK7TYACAAbAAAAMhkFCBM4IgHdkvgHQKiYAAAAAGQEAg0RrDs2AAkAGwAAADIZBQgTUycB3ZP/B0CpzQAHAABkBAINERQHNgAKABsAAAAyGQUIFAgpAd2WEgdAqZ4AAwAAZAQCDw/suTYACwAbAAAAMhkFCBQjNAHdlEYHQKerAAUAAGQEAg4O5Yk2AAwAGwAAADIZBQgUODgB3ZQjB0Cl7gAHAABkBAIMDrxsNgANABsAAAAyGQUIFFNDAd2WXAdAp5MAEAAAZAQCDQ7IQTYADgAbAAAAMhkFCBUQIgHdlDcHQKkeAAAAAGQEAg0Oo1WjAg==`)
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

		service.GenGpsEventFromSinglePacket(singlePacket, "seqNo", "cntrNum", "revTime")
	}
}

func TestSinglePacketSth(t *testing.T) {
	decoded, err := base64.StdEncoding.DecodeString(`NwAPAZg2AAMAGwAAADIZBQgSIwQB3ZPdB0CqhwAAAABkBAIOEcq3NgAEABsAAAAyGQUIEjgQAd2UQgdAqVsAAQAAZAQCDhE71jYABQAbAAAAMhkFCBJTEwHdk6oHQKmhAAAAAGQEAg0QeBo2AAYAGwAAADIZBQgTCBcB3ZOSB0CopwAAAABkBAIOEWwINgAHABsAAAAyGQUIEyMhAd2VWwdApgYAAwAAZAQCEBHK7TYACAAbAAAAMhkFCBM4IgHdkvgHQKiYAAAAAGQEAg0RrDs2AAkAGwAAADIZBQgTUycB3ZP/B0CpzQAHAABkBAINERQHNgAKABsAAAAyGQUIFAgpAd2WEgdAqZ4AAwAAZAQCDw/suTYACwAbAAAAMhkFCBQjNAHdlEYHQKerAAUAAGQEAg4O5Yk2AAwAGwAAADIZBQgUODgB3ZQjB0Cl7gAHAABkBAIMDrxsNgANABsAAAAyGQUIFFNDAd2WXAdAp5MAEAAAZAQCDQ7IQTYADgAbAAAAMhkFCBUQIgHdlDcHQKkeAAAAAGQEAg0Oo1WjAg==`)
	if err != nil {
		return
	}
	_, reqDataSegment := parseBytes(decoded)
	singlePacket, err := service.ParseToSinglePacket(reqDataSegment)
	if err != nil {
		return
	}

	service.GenGpsEventFromSinglePacket(singlePacket, "seqNo", "cntrNum", "revTime")
}
