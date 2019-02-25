package util

import (
	"math"
	"reflect"

	entity "github.com/leeshing0315/flogo-components/common/entity"
	crg "github.com/leeshing0315/go-city-reverse-geocoder"
)

const EARTH_RADIAS = 6378137
const DISTANCE_FROM_CITY = 10

var geofences = make(map[string][]Geofence)

func isPointInPolygon(point [2]float64, polyCorners [][2]float64) bool {
	var i int
	var j = len(polyCorners) - 1
	var oddNodes = false
	for i = 0; i < len(polyCorners); i++ {
		if polyCorners[i][1] < point[1] && polyCorners[j][1] >= point[1] ||
			polyCorners[j][1] < point[1] && polyCorners[i][1] >= point[1] {
			if polyCorners[i][0]+(point[1]-polyCorners[i][1])/(polyCorners[j][1]-polyCorners[i][1])*(polyCorners[j][0]-polyCorners[i][0]) < point[0] {
				oddNodes = !oddNodes
			}
		}
		j = i
	}
	return oddNodes
}

func calculateDistance(point1 [2]float64, point2 [2]float64) float64 {
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
		if geo.geoType == "circle" {
			var tmp [2]float64 = geo.coords.coordinates.([2]float64)
			distance := calculateDistance(tmp, [2]float64{lon, lat})
			if distance <= geo.radiusInMetre.(float64) {
				result = geo
				return result
			}
		} else {
			if geo.coords.coordinates != nil {
				var inPolygon = isPointInPolygon([2]float64{lon, lat}, geo.coords.coordinates.([][][2]float64)[0])
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
		if reflect.ValueOf(oceanResult).IsValid() {
			result = oceanResult
		} else {
			result = location[0]
		}
	} else {
		result = location[0]
	}
	return result
}

func searchFromOceanPolygon(lat float64, lon float64) Geofence {
	var result Geofence
	for i := 0; i < len(geofences["ocean"]); i++ {
		geo := geofences["ocean"][i]
		if geo.coords.coordinates != nil {
			var inPolygon = isPointInPolygon([2]float64{lon, lat}, geo.coords.coordinates.([][][2]float64)[0])
			if inPolygon {
				result = geo
				return result
			}
		}
	}
	return result
}

func AttachLocation (lat float64, lon float64, gpsevent entity.GpsEvent) entity.GpsEvent {
	var location interface{}
	if (gpsevent.Carrier == "COSU") {
		location = getLocationByLatLon(lat, lon, "COSCO");
	} else {
		location = getLocationByLatLon(lat, lon, gpsevent.Carrier.(string));
	}
	if (location != nil) {
		//typeName := reflect.TypeOf(location).String()
		_, ok := location.(crg.Result)
		if (ok) {
			tmpLocation := location.(crg.Result)
			gpsevent.Address = entity.GpsEventAddress{
				Distance: tmpLocation.Distance,
				City: tmpLocation.City,
				RegionCode: tmpLocation.Region_code,
				Region: tmpLocation.Region,
				CountryCode: tmpLocation.Country_code,
				Country: tmpLocation.Country,
			}
			if (tmpLocation.City != "" && tmpLocation.Region != "" && tmpLocation.Country != "") {
				gpsevent.DisplayName = tmpLocation.City + ", " + tmpLocation.Region + ", " + tmpLocation.Country
			} else {
				gpsevent.DisplayName = ""
			}
		}
		_, ok = location.(Geofence)
		if (ok) {
			tmpLocation := location.(Geofence)
			gpsevent.Address = entity.GpsEventAddress{city: tmpLocation.geoCity, country: tmpLocation.geoCountry, name: tmpLocation.geoName, code: tmpLocation.geoCode}
			gpsevent.DisplayName = tmpLocation.geoName
			if (tmpLocation.geoCity != "") {
				gpsevent.DisplayName = gpsevent.DisplayName.string() + ", " + tmpLocation.geoCity
			}
			if (tmpLocation.geoCountry != "") {
				gpsevent.DisplayName = gpsevent.DisplayName.string() + ", " + tmpLocation.geoCountry
			}
		}
		return gpsevent
	} else {
		gpsevent.DisplayName = ""
		return gpsevent
	}
}

//func main() {
//	var coord = Coords{
//		Type: "Polygon",
//		coordinates: [][][2]float64{
//			{
//				{113.078, 22.905},
//				{113.079, 22.903},
//				{113.077, 22.902},
//				{113.078, 22.903},
//				{113.078, 22.905},
//			},
//		},
//	}
//	var coord2 = Coords{
//		Type: "Polygon",
//		coordinates: [][][2]float64{
//			{
//				{113.074, 22.905},
//				{113.073, 22.903},
//				{113.076, 22.902},
//				{113.077, 22.903},
//				{113.074, 22.905},
//			},
//		},
//	}
//	var tmp1 = Geofence{
//		geoId:   "bbd9030d-18a1-4d49-b349-44337d70bb22",
//		geoName: "Kerry Intermodal Services",
//		geoType: "polygon",
//		coords: coord,
//		isDeleted:  "F",
//		createdAt:  "2017-03-02T09:49:18.007Z",
//		geoCode:    "ADL51",
//		geoLocType: "Rail Ramp",
//		isDisabled: false,
//	};
//	var tmp2 = Geofence{
//		geoId:   "bbd9030d-18a1-4d49-b349-44337d70bb22",
//		geoName: "testestestest",
//		geoType: "polygon",
//		coords: coord2,
//		isDeleted:  "F",
//		createdAt:  "2017-03-02T09:49:18.007Z",
//		geoCode:    "ADL51",
//		geoLocType: "Rail Ramp",
//		isDisabled: true,
//	};
//	geofences["COSCO"] = []Geofence{tmp1}
//	geofences["ocean"] = []Geofence{tmp2}
//	//print("asdfasdfsadfasdfasdfs")
//	//print(geofences["COSCO"])
//	var result = getLocationByLatLon(22.904, 113.074, "COSCO")
//	fmt.Println(result, "result");
//	print(reflect.TypeOf(result).String() == "[]geocoder.Result")
//	//print(reflect.ValueOf(result).String())
//	//geofences2 :=[][2]float64{{113.074, 22.905},{113.073, 22.903},{113.076, 22.902},{113.077, 22.903},{113.074, 22.905}}
//	//var result = isPointInPolygon ([2]float64{113.074, 22.904}, geofences2)
//	//print(result)
//	//t.is(result, 1113.1949)
//}

type Coords struct {
	Type        string `json:"type"`
	coordinates interface{}
}

type Geofence struct {
	geoId         interface{}
	geoName       interface{}
	coords        Coords
	isDeleted     interface{}
	createdAt     interface{}
	geoCode       interface{}
	geoLocType    interface{}
	isDisabled    bool
	radiusInMetre interface{}
	geoCity       interface{}
	source        interface{}
	carrier       interface{}
	geoCountry    interface{}
	geoType       interface{}
}

type Location struct {
	country      interface{}
	country_code interface{}
	region       interface{}
	region_code  interface{}
	city         interface{}
	latitude     float64
	longitude    float64
	distance     float64
}
