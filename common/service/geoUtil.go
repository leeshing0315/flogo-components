package service

import (
	"context"
	"os"

	"math"
	"reflect"
	"strconv"

	"github.com/leeshing0315/flogo-components/common/entity"
	crg "github.com/leeshing0315/go-city-reverse-geocoder"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const EARTH_RADIAS = 6378137
const DISTANCE_FROM_CITY = 10

var geofences = make(map[string][]map[string]interface{})

var mongoUri = os.Getenv("MONGO_URI")

func init() {
	if mongoUri == "" {
		return
	}
	option := options.Client()
	option.ApplyURI(mongoUri)
	client, err := mongo.Connect(context.Background(), option)
	if err != nil {
		return
	}
	iotDatabase := client.Database("iot")
	geoColl := iotDatabase.Collection("geofences")
	coscoGeo, err := queryGeofences(geoColl, map[string]interface{}{
		"carrier":   "COSCO",
		"isDeleted": "F",
	})
	if err != nil {
		return
	}
	ooclGeo, err := queryGeofences(geoColl, map[string]interface{}{
		"carrier":   "OOCL",
		"isDeleted": "F",
	})
	if err != nil {
		return
	}
	oceanGeo, err := queryGeofences(geoColl, map[string]interface{}{
		"carrier": map[string]interface{}{
			"$exists": false,
		},
		"geoLocType": "Ocean",
		"isDeleted":  "F",
	})
	if err != nil {
		return
	}
	geofences["COSCO"] = coscoGeo
	geofences["OOCL"] = ooclGeo
	geofences["OCEAN"] = oceanGeo
}

func queryGeofences(geoColl *mongo.Collection, condition map[string]interface{}) ([]map[string]interface{}, error) {
	cur, err := geoColl.Find(context.Background(), condition)
	if err != nil {
		return nil, err
	}
	var result []map[string]interface{}
	for cur.Next(context.Background()) {
		var elem map[string]interface{}
		// elem := Geofence{}
		if decodeErr := cur.Decode(&elem); decodeErr != nil {
			return nil, err
		}
		result = append(result, elem)
	}
	return result, nil
}

func isPointInPolygon(point []float64, polyCornersTmp []interface{}) bool {
	var i int
	polyCorners := []interface{}(polyCornersTmp[0].(primitive.A))
	var j = len(polyCorners) - 1
	var oddNodes = false
	for i = 0; i < len(polyCorners); i++ {
		polyCornersItemCurrent := []interface{}(polyCorners[i].(primitive.A))
		polyCornersItemBefore := []interface{}(polyCorners[j].(primitive.A))
		polyCornersItemCurrentArr := []float64{}
		polyCornersItemBeforeArr := []float64{}
		if reflect.TypeOf(polyCornersItemBefore[0]).Name() == "int32" {
			polyCornersItemBeforeArr = append(polyCornersItemBeforeArr, float64(polyCornersItemBefore[0].(int32)))
		}
		if reflect.TypeOf(polyCornersItemCurrent[0]).Name() == "float64" {
			polyCornersItemCurrentArr = append(polyCornersItemCurrentArr, polyCornersItemCurrent[0].(float64))
		}
		if reflect.TypeOf(polyCornersItemBefore[1]).Name() == "int32" {
			polyCornersItemBeforeArr = append(polyCornersItemBeforeArr, float64(polyCornersItemBefore[1].(int32)))
		}
		if reflect.TypeOf(polyCornersItemCurrent[1]).Name() == "float64" {
			polyCornersItemCurrentArr = append(polyCornersItemCurrentArr, polyCornersItemCurrent[1].(float64))
		}
		if reflect.TypeOf(polyCornersItemBefore[0]).Name() == "float64" {
			polyCornersItemBeforeArr = append(polyCornersItemBeforeArr, polyCornersItemBefore[0].(float64))
		}
		if reflect.TypeOf(polyCornersItemCurrent[0]).Name() == "int32" {
			polyCornersItemCurrentArr = append(polyCornersItemCurrentArr, float64(polyCornersItemCurrent[0].(int32)))
		}
		if reflect.TypeOf(polyCornersItemBefore[1]).Name() == "float64" {
			polyCornersItemBeforeArr = append(polyCornersItemBeforeArr, polyCornersItemBefore[1].(float64))
		}
		if reflect.TypeOf(polyCornersItemCurrent[1]).Name() == "int32" {
			polyCornersItemCurrentArr = append(polyCornersItemCurrentArr, float64(polyCornersItemCurrent[1].(int32)))
		}
		if polyCornersItemCurrentArr[1] < point[1] && polyCornersItemBeforeArr[1] >= point[1] ||
			polyCornersItemBeforeArr[1] < point[1] && polyCornersItemCurrentArr[1] >= point[1] {
			if polyCornersItemCurrentArr[0]+(point[1]-polyCornersItemCurrentArr[1])/(polyCornersItemBeforeArr[1]-polyCornersItemCurrentArr[1])*(polyCornersItemBeforeArr[0]-polyCornersItemCurrentArr[0]) < point[0] {
				oddNodes = !oddNodes
			}
		}
		j = i
	}
	return oddNodes
}

func calculateDistance(point1 []float64, point2 []float64) float64 {
	var radLat1 = point1[0] * math.Pi / 180.0
	var radLat2 = point2[0] * math.Pi / 180.0
	var a = radLat1 - radLat2
	var b = point1[1]*math.Pi/180.0 - point2[1]*math.Pi/180.0
	var s = 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	s = s * EARTH_RADIAS
	var tmp float64
	if s*10000 > 0 {
		tmp = math.Trunc(s*10000 + 0.5)
	} else {
		tmp = math.Trunc(s*10000 - 0.5)
	}
	s = tmp / 10000
	return s
}

func getLocationByLatLon(lat float64, lon float64, carrier string) interface{} {
	if geofences == nil || len(geofences[carrier]) == 0 {
		return nil
	}
	//if (lat == "" || lon == "") {
	//	return nil;
	//}
	var result interface{}
	for i := 0; i < len(geofences[carrier]); i++ {
		geo := geofences[carrier][i]
		if geo["geoType"] == "circle" {
			var tmp = []interface{}((geo["coords"]).(map[string]interface{})["coordinates"].(primitive.A))
			var point []float64
			for j := 0; j < len(tmp); j++ {
				point = append(point, tmp[j].(float64))
			}
			distance := calculateDistance(point, []float64{lon, lat})
			radiusInMetre, ok := geo["radiusInMetre"].(float64)
			if !ok {
				continue
			}
			if distance <= radiusInMetre {
				result = geo
				return result
			}
		}
		if geo["geoType"] == "polygon" {
			coordinates := geo["coords"].(map[string]interface{})["coordinates"]
			if coordinates != nil {
				var inPolygon = isPointInPolygon([]float64{lon, lat}, []interface{}(coordinates.(primitive.A)))
				if inPolygon {
					result = geo
					return result
				}
			}
		}
	}
	var location, error = crg.GetNearestCities(lat, lon, 1, "mi")
	if location[0].Distance > DISTANCE_FROM_CITY || error != nil {
		var oceanResult = searchFromOceanPolygon(lat, lon)
		if reflect.ValueOf(oceanResult).IsValid() && oceanResult["geoName"] != "" {
			result = oceanResult
		} else {
			result = location[0]
		}
	} else {
		result = location[0]
	}
	return result
}

func searchFromOceanPolygon(lat float64, lon float64) map[string]interface{} {
	var result map[string]interface{}
	for i := 0; i < len(geofences["OCEAN"]); i++ {
		geo := geofences["OCEAN"][i]
		coordinates := (geo["coords"]).(map[string]interface{})["coordinates"]
		if coordinates != nil {
			var inPolygon = isPointInPolygon([]float64{lon, lat}, []interface{}(coordinates.(primitive.A)))
			if inPolygon {
				result = geo
				return result
			}
		}
	}
	return result
}

func AttachLocation(gpsevent *entity.GpsEvent) *entity.GpsEvent {
	var location interface{}
	lat, err := strconv.ParseFloat(gpsevent.Lat, 64)
	if err != nil {
		return nil
	}
	lon, err := strconv.ParseFloat(gpsevent.Lng, 64)
	if err != nil {
		return nil
	}
	location = getLocationByLatLon(lat, lon, "COSCO")
	if location != nil {
		_, ok := location.(crg.Result)
		if ok {
			tmpLocation := location.(crg.Result)
			gpsevent.Address = entity.GpsEventAddress{
				Distance:    tmpLocation.Distance,
				City:        tmpLocation.City,
				RegionCode:  tmpLocation.Region_code,
				Region:      tmpLocation.Region,
				CountryCode: tmpLocation.Country_code,
				Country:     tmpLocation.Country,
			}
			if tmpLocation.City != "" && tmpLocation.Region != "" && tmpLocation.Country != "" {
				gpsevent.DisplayName = tmpLocation.City + ", " + tmpLocation.Region + ", " + tmpLocation.Country
			} else {
				gpsevent.DisplayName = ""
			}
		}
		_, ok = location.(map[string]interface{})
		if ok {
			tmpLocation := location.(map[string]interface{})
			gpsevent.Address = entity.GpsEventAddress{}
			if tmpLocation["geoName"] != nil {
				gpsevent.Address.Name = tmpLocation["geoName"].(string)
			}
			if tmpLocation["geoCity"] != nil {
				gpsevent.Address.City = tmpLocation["geoCity"].(string)
			}
			if tmpLocation["geoCountry"] != nil {
				gpsevent.Address.Country = tmpLocation["geoCountry"].(string)
			}
			if tmpLocation["geoCode"] != nil {
				gpsevent.Address.Code = tmpLocation["geoCode"].(string)
			}
			gpsevent.DisplayName = tmpLocation["geoName"].(string)
			if tmpLocation["geoCity"] != nil {
				gpsevent.DisplayName = gpsevent.DisplayName + ", " + tmpLocation["geoCity"].(string)
			}
			if tmpLocation["geoCountry"] != nil {
				gpsevent.DisplayName = gpsevent.DisplayName + ", " + tmpLocation["geoCountry"].(string)
			}
		}
	} else {
		gpsevent.DisplayName = ""
	}
	ooclLocation := getLocationByLatLon(lat, lon, "OOCL")
	if ooclLocation != nil {
		_, ok := ooclLocation.(map[string]interface{})
		if ok {
			tmpLocation := location.(map[string]interface{})
			gpsevent.Address.OoclDisplayName = ""
			if tmpLocation["geoName"] != nil {
				gpsevent.Address.OoclName = tmpLocation["geoName"].(string)
				gpsevent.Address.OoclDisplayName = tmpLocation["geoName"].(string)
			}
			if tmpLocation["geoCode"] != nil {
				gpsevent.Address.OoclCode = tmpLocation["geoCode"].(string)
			}
			if tmpLocation["geoCity"] != nil {
				gpsevent.Address.OoclDisplayName = gpsevent.Address.OoclDisplayName + ", " + tmpLocation["geoCity"].(string)
			}
			if tmpLocation["geoCountry"] != nil {
				gpsevent.Address.OoclDisplayName = gpsevent.Address.OoclDisplayName + ", " + tmpLocation["geoCountry"].(string)
			}
		}
	}

	return gpsevent
}
