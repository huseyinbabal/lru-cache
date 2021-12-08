package lru

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

var cache = New(255)

type IpResponse struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func getIpGeoLocationWithCache(ip string) IpResponse {
	val := cache.Get([]byte(ip))
	if val != nil {
		log.Println("Getting from cache " + ip)
		var data IpResponse
		json.Unmarshal(val, &data)
		return data
	}
	ipData := getIpGeoLocation(ip)
	ipDataBytes := new(bytes.Buffer)
	json.NewEncoder(ipDataBytes).Encode(ipData)
	log.Printf("Puting into cache %v", ipData)
	cache.Put([]byte(ip), ipDataBytes.Bytes())
	return ipData
}

func getIpGeoLocation(ip string) IpResponse {
	resp, err := http.Get("http://ip-api.com/json/" + ip)
	if err != nil {
		log.Fatalf("http get err %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read err %v", err)
	}

	var data IpResponse
	json.Unmarshal(body, &data)
	return data
}
func BenchmarkGetWithoutCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getIpGeoLocation(fmt.Sprintf("24.48.9.%d", i))
	}
}

func BenchmarkGetWithCache(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getIpGeoLocationWithCache(fmt.Sprintf("24.48.9.%d", i))
	}
}
