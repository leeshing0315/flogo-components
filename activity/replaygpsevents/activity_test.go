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
	decoded, err := base64.StdEncoding.DecodeString(`NwDqA2k2AOEAWgCAAHIZBRUJNCQCHwa3BzbmnwEmAIZkBAIQEgIsRAH////DwzfDI8NX7D3///8nJhDXBOkC/sSPAHfC6wJD/+0AKgBqAH4Ag/8IAkRdCQj96v3q/er96gUBw/GONgDiAFoAgAByGQUVCTcTAh7kEgc3EHIBJgCGZAQCEBICLEQB////wsM3w0vDV+xG////JwcQ7gB0AADEdwB8wyMBUAB2ABwAAACDAIP/CAIAJAkI/er96v3q/eoFAcL/pzYA4wBaAIAAchkFFQk4NwIe0vEHNyV3ASgAh2QEAhASAixEAf///8PDN8Mjw1fsSP///ycPENoE3QLtxIsAesL3Ajz/+QAoAGoAfwCD/wgCRF0JCP3q/er96v3qBQHDE8Y2AOQAWgCAAHIZBRUJQScCHrALBzdQzAEoAIZkBAIPEgIsRAH////CwzfDS8NX7DD///8nKRDsAHMAAMRrAHvDNwFUAHgAHQAAAIIAg/8IAgAkCQj96v3q/er96gUBwg+kNgDlAFoAgAByGQUVCUQRAh6OEgc3ecMBKACGZAQCDhICLEQB////w8M3w0PDV+w3////JygQ2QTiAvXEjwB6ww8COf/2ACgAagB/AIP/CAJEXQkI/er96v3q/eoFAcO6WjYA5gBaAIAAchkFFQlGFwIedFYHN5nWASYAg2QEAg4SAixEAf///8LDN8NLw1fsRv///ycHEOwAdAAAxGsAe8MjAUAAdgAdAAAAgwCD/wgCACQJCP3q/er96v3qBQHC3Kc2AOcAWgCAAHIZBRUJSAICHmBuBze2JAEmAIFkBAIOEgIsRAH////DwzfDI8NX7Dz///8nKBDXBNkC6sSPAHfC2wI+/+8AKwBqAH0Ag/8IAkRdCQj96v3q/er96gUBw1KqNgDoAFoAgAByGQUVCVBQAh5ByAc35EsBKACBZAQCDhICLEQB////wsM3w0vDV+xH////JwcQ6wB0AADEdwB7wysBQQB4AB0AAACDAIP/CAIAJAkI/er96v3q/eoFAcKhojYA6QBaAIAAchkFFQlSNQIeLj4HOAFVASoAgWQEAg4SAixEAf///8PDN8Mjw1fsPP///ycoENkE0gLhxJcAd8LbAj3/7gArAGoAfQCD/wgCRF0JCP3q/er96v3qBQHDPMkfdQ==`)
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
	decoded, err := base64.StdEncoding.DecodeString(`NgDJAFoAgAByGQUVCRFTAh/oRQc1jwoBDgCEZAQCEBICLEQB////w8M3wyPDV+w+////JyYQ1wTeAu3EiwB5wusCS//tACoAagB+AIP/CAJEXQkI/er96v3q/eoFAcP9kw==`)
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
