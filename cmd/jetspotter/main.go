package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// Structs
type FlightData struct {
	AC    []Aircraft `json:"ac"`
	Msg   string     `json:"msg"`
	Now   int64      `json:"now"`
	Total int        `json:"total"`
	Ctime int64      `json:"ctime"`
	Ptime int        `json:"ptime"`
}

type Aircraft struct {
	Hex         string        `json:"hex"`
	Type        string        `json:"type"`
	Flight      string        `json:"flight"`
	R           string        `json:"r"`
	T           string        `json:"t"`
	Desc        string        `json:"desc"`
	AltBaro     interface{}   `json:"alt_baro"`
	AltGeom     int           `json:"alt_geom"`
	GS          float64       `json:"gs"`
	IAS         int           `json:"ias"`
	TAS         int           `json:"tas"`
	Mach        float64       `json:"mach"`
	WD          int           `json:"wd"`
	WS          int           `json:"ws"`
	OAT         int           `json:"oat"`
	TAT         int           `json:"tat"`
	Track       float64       `json:"track"`
	TrackRate   float64       `json:"track_rate"`
	Roll        float64       `json:"roll"`
	MagHeading  float64       `json:"mag_heading"`
	TrueHeading float64       `json:"true_heading"`
	BaroRate    int           `json:"baro_rate"`
	GeomRate    int           `json:"geom_rate"`
	Squawk      string        `json:"squawk"`
	Emergency   string        `json:"emergency"`
	Category    string        `json:"category"`
	NavQNH      float64       `json:"nav_qnh"`
	NavAltMCP   int           `json:"nav_altitude_mcp"`
	NavAltFMS   int           `json:"nav_altitude_fms"`
	NavHeading  float64       `json:"nav_heading"`
	Lat         float64       `json:"lat"`
	Lon         float64       `json:"lon"`
	NIC         int           `json:"nic"`
	RC          int           `json:"rc"`
	SeenPos     float64       `json:"seen_pos"`
	Version     int           `json:"version"`
	NICBaro     int           `json:"nic_baro"`
	NACP        int           `json:"nac_p"`
	NACV        int           `json:"nac_v"`
	SIL         int           `json:"sil"`
	SILType     string        `json:"sil_type"`
	GVA         int           `json:"gva"`
	SDA         int           `json:"sda"`
	Alert       int           `json:"alert"`
	SPI         int           `json:"spi"`
	MLAT        []interface{} `json:"mlat"`
	TISB        []interface{} `json:"tisb"`
	Messages    int           `json:"messages"`
	Seen        float64       `json:"seen"`
	RSSI        float64       `json:"rssi"`
	Dst         float64       `json:"dst"`
	Dir         float64       `json:"dir"`
}

// Vars
var (
	baseURL = "https://api.adsb.one/v2/point"
)

func exitWithError(err error) {
	fmt.Printf("Something went wrong: %v\n", err)
	os.Exit(1)
}

/* getAircraftInProximity returns all aircraft within a specified maxRange of a latitude/longitude point. */
func getAircraftInProximity(latitude string, longitude string, maxRange int) (aircraft []Aircraft, err error) {
	var flightData FlightData
	endpoint, err := url.JoinPath(baseURL, latitude, longitude, strconv.Itoa(maxRange))
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &flightData)
	if err != nil {
		return nil, err
	}

	return flightData.AC, nil
}

func main() {

	var vipers []Aircraft
	aircraft, err := getAircraftInProximity("51.078395", "5.018769", 30)
	if err != nil {
		exitWithError(err)
	}

	for _, ac := range aircraft {
		if ac.Flight == "" {
			ac.Flight = "UNKNOWN"
		}

		if ac.T == "F16" {
			vipers = append(vipers, ac)
		}

		fmt.Printf("=== %s ===\n", ac.Flight)
		fmt.Printf("Description: %s\nType: %s\nTail number: %s\nAltitude: %v\n", ac.Desc, ac.T, ac.R, ac.AltBaro)
		fmt.Println()
	}

	if len(vipers) > 0 {
		fmt.Println("****** SPOTTED AN F16, SHOULD SEND SLACK MESSAGE! ******")
		for _, viper := range vipers {
			// Calculate distance based on lat lon
			fmt.Printf("CALLSIGN: %s\nSQUAWK: %s\nTAIL: %s\nLAT: %v\nLON: %v\n\n", viper.Flight, viper.Squawk, viper.R, viper.Lat, viper.Lon)

		}
	}
}
