package smucntrdevmapping

import (
	"encoding/json"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/leeshing0315/flogo-components/common/entity"
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
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {
	cntrDevMappingStr := context.GetInput("cntrDevMapping").(string)
	devtype := context.GetInput("devtype").(string)
	ip := context.GetInput("ip").(string)
	firmwareVersion := context.GetInput("firmwareVersion").(string)

	// do eval
	var cntrDevMapping = make(map[string]interface{})
	err = json.Unmarshal([]byte(cntrDevMappingStr), &cntrDevMapping)
	if err != nil {
		return false, err
	}

	deviceInfo := &entity.DeviceInfo{}
	if cntrDevMapping["simno"] != nil {
		deviceInfo.Simno = cntrDevMapping["simno"].(string)
	}
	deviceInfo.Devtype = devtype
	deviceInfo.Ip = ip
	deviceInfo.Remark = firmwareVersion
	deviceInfo.Setaddr = "1"
	deviceInfo.TableName = "Tbldevinfo"
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	deviceInfo.Savetime = time.Now().In(loc).Format("2006-01-02 15:04:05.0")
	deviceInfoBytes, err := json.Marshal(deviceInfo)
	if err == nil {
		context.SetOutput("deviceInfo", string(deviceInfoBytes))
	}

	context.SetOutput("sim", cntrDevMapping["simno"])
	context.SetOutput("devId", cntrDevMapping["carno"])
	context.SetOutput("cntrNum", cntrDevMapping["carid"])
	if cntrDevMapping["company"] == nil || cntrDevMapping["company"].(string) == "" {
		context.SetOutput("company", "COSU")
	} else if cntrDevMapping["company"].(string) == "OOCL" {
		context.SetOutput("company", "OOLU")
	} else {
		context.SetOutput("company", cntrDevMapping["company"].(string))
	}

	println("**********login*sim", cntrDevMapping["simno"].(string), "**********")
	println("**********login*devId", cntrDevMapping["carno"].(string), "**********")
	println("**********login*cntrNum", cntrDevMapping["carid"].(string), "**********")

	return true, nil
}
