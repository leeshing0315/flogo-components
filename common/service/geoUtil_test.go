package service

import (
	"fmt"
	"testing"

	"github.com/leeshing0315/flogo-components/common/entity"
)

func TestGeoUtil(t *testing.T) {
	// [6]:195
	// [7]:87
	// [8]:254
	// [9]:143
	// [10]:254
	// [11]:131
	gpsEvents := entity.GpsEvent{
		CntrNum: "test",
		Lng:     "128",
		Lat:     "20",
		Carrier: "COSU",
	}
	// test()
	gps := AttachLocation(&gpsEvents)
	// result := getLocationByLatLon(22.904, 113.074, "COSCO")
	fmt.Println(gps)
}
