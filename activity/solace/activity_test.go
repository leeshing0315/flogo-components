package solace

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/stretchr/testify/assert"
)

const GPS_EVENT_SAMPLE = "{ " +
	"\"seqno\" : \"11129352\", " +
	"\"cntrNum\" : \"CXRU1035624\", " +
	"\"revTime\" : \"2018-07-19T00:00:48+08:00\", " +
	"\"cltTime\" : \"2018-07-19T00:00:31+08:00\", " +
	"\"locateTime\" : \"2018-07-19T00:00:31+08:00\", " +
	"\"eleState\" : \"1\", " +
	"\"batLevel\" : \"100\", " +
	"\"opMode\" : \"Modulation\", " +
	"\"setTem\" : \"0\", " +
	"\"supTem\" : \"0.2\", " +
	"\"retTem\" : \"2.3\", " +
	"\"lng\" : \"139.759088\", " +
	"\"lat\" : \"35.624676\", " +
	"\"speed\" : \"0\", " +
	"\"direction\" : \"340\", " +
	"\"hpt\" : \"1090\", " +
	"\"posFlag\" : \"1\", " +
	"\"ism\" : \"0\", " +
	"\"gpsNum\" : \"12\", " +
	"\"bdNum\" : \"12\", " +
	"\"lpt\" : \"30\", " +
	"\"pt\" : \"409.6\", " +
	"\"ct1\" : \"13.1\", " +
	"\"ct2\" : \"9.2\", " +
	"\"ambs\" : \"29.9\", " +
	"\"eis\" : \"-0.7\", " +
	"\"eos\" : \"-0.5\", " +
	"\"dchs\" : \"78.1\", " +
	"\"sgs\" : \"-19.5\", " +
	"\"smv\" : \"13.1\", " +
	"\"ev\" : \"81\", " +
	"\"dss\" : \"0\", " +
	"\"drs\" : \"2.1\", " +
	"\"source\" : \"WEB_SERVICE\", " +
	"\"carrier\" : \"COSU\", " +
	"\"displayName\" : \"Shinagawa, Tokyo, Japan\", " +
	"\"address\" : {" +
	"\"distance\" : 1.5332328884126547, " +
	"\"longitude\" : 139.748, " +
	"\"latitude\" : 35.6044, " +
	"\"city\" : \"Shinagawa\", " +
	"\"region_code\" : \"13\", " +
	"\"region\" : \"Tokyo\", " +
	"\"country_code\" : \"JP\", " +
	"\"country\" : \"Japan\"" +
	"}, " +
	"}"

const (
	TEST_HOSTIP      = "localhost:5672"
	TEST_VPN         = "default"
	TEST_TOPICNAME   = "TOPIC/TEST"
	TEST_DATA        = GPS_EVENT_SAMPLE
	TEST_ROUTINE_NUM = 10000
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
	tc.SetInput("hostIP", TEST_HOSTIP)
	tc.SetInput("vpnName", TEST_VPN)
	tc.SetInput("topicName", TEST_TOPICNAME)
	tc.SetInput("data", TEST_DATA)
	ch := make(chan string)
	for i := 0; i < TEST_ROUTINE_NUM; i++ {
		go concurrent(ch, act, tc, t)
	}
	for i := 0; i < TEST_ROUTINE_NUM; i++ {
		<-ch
	}
}

func concurrent(ch chan string, act activity.Activity, tc *test.TestActivityContext, t *testing.T) {
	act.Eval(tc)
	result := tc.GetOutput("publishSuccess")
	assert.Equal(t, "true", result)
	ch <- result.(string)
}
