package jetspotter

// Structs
/* FlightData is a struct of the json received by the ADS-B api. */
type FlightData struct {
	AC    []Aircraft `json:"ac"`
	Msg   string     `json:"msg"`
	Now   int64      `json:"now"`
	Total int        `json:"total"`
	Ctime int64      `json:"ctime"`
	Ptime int        `json:"ptime"`
}

type Aircraft struct {
	ICAO        string        `json:"hex"`
	Type        string        `json:"type"`
	Callsign    string        `json:"flight"`
	TailNumber  string        `json:"r"`
	PlaneType   string        `json:"t"`
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
