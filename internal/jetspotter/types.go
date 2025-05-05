package jetspotter

// FlightData is a struct of the json received by the ADS-B api
type FlightData struct {
	// A slice of aircrafts
	AC    []AircraftRaw `json:"ac"`
	Msg   string        `json:"msg"`
	Now   int64         `json:"now"`
	Total int           `json:"total"`
	Ctime int64         `json:"ctime"`
	Ptime int           `json:"ptime"`
}

// AircraftRaw contains all the metadata of an aircraft as defined by the ADS-B
// https://www.adsbexchange.com/ads-b-data-field-explanations/
type AircraftRaw struct {
	// Aircraft Type Designator number, basically the unique identifier of an aircraft
	ICAO string `json:"hex"`
	// Type of transponder used that received the data
	Type string `json:"type"`
	// Callsign or flight name of the aircraft, if not set 'NONE' is used
	Callsign string `json:"flight"`
	// Tail number of the aircraft
	Registration string `json:"r"`
	// Type of the aircraft
	PlaneType string `json:"t"`
	Desc      string `json:"desc"`
	// Barometric altitude in feet
	AltBaro interface{} `json:"alt_baro"`
	// Geometric (GNSS / INS) altitude in feet referenced to the WGS84 ellipsoid
	AltGeom int `json:"alt_geom"`
	// Ground speed in knots
	GS float64 `json:"gs"`
	// Indicated air speed in knots
	IAS int `json:"ias"`
	// True air speed in knots
	TAS int `json:"tas"`
	// Mach number
	Mach      float64 `json:"mach"`
	WD        int     `json:"wd"`
	WS        int     `json:"ws"`
	OAT       int     `json:"oat"`
	TAT       int     `json:"tat"`
	Track     float64 `json:"track"`
	TrackRate float64 `json:"track_rate"`
	Roll      float64 `json:"roll"`
	// Heading, degrees clockwise from magnetic north
	MagHeading  float64 `json:"mag_heading"`
	TrueHeading float64 `json:"true_heading"`
	BaroRate    int     `json:"baro_rate"`
	GeomRate    int     `json:"geom_rate"`
	// Mode A code (Squawk), encoded as 4 octal digits
	Squawk     string  `json:"squawk"`
	Emergency  string  `json:"emergency"`
	Category   string  `json:"category"`
	NavQNH     float64 `json:"nav_qnh"`
	NavAltMCP  int     `json:"nav_altitude_mcp"`
	NavAltFMS  int     `json:"nav_altitude_fms"`
	NavHeading float64 `json:"nav_heading"`
	// Aircraft latitude position in decimal degrees
	Lat float64 `json:"lat"`
	// Aircraft longitude position in decimal degrees
	Lon float64 `json:"lon"`
	// Database flags, 1 = military, 2 = interesting, 3 = both
	DbFlags  int           `json:"dbFlags"`
	NIC      int           `json:"nic"`
	RC       int           `json:"rc"`
	SeenPos  float64       `json:"seen_pos"`
	Version  int           `json:"version"`
	NICBaro  int           `json:"nic_baro"`
	NACP     int           `json:"nac_p"`
	NACV     int           `json:"nac_v"`
	SIL      int           `json:"sil"`
	SILType  string        `json:"sil_type"`
	GVA      int           `json:"gva"`
	SDA      int           `json:"sda"`
	Alert    int           `json:"alert"`
	SPI      int           `json:"spi"`
	MLAT     []interface{} `json:"mlat"`
	TISB     []interface{} `json:"tisb"`
	Messages int           `json:"messages"`
	Seen     float64       `json:"seen"`
	RSSI     float64       `json:"rssi"`
	Dst      float64       `json:"dst"`
	Dir      float64       `json:"dir"`
}

// Aircraft contains all fields that we want to print, regardless of which medium is used
type Aircraft struct {
	// Aircraft Type Designator number, basically the unique identifier of an aircraft
	ICAO string

	// Callsign or flight name of the aircraft, if not set 'NONE' is used
	Callsign string

	// Description of the aircraft
	Description string

	// Type of the aircraft
	Type string

	// Latitude of the aircraft in decimal degrees
	Latitude float64

	// Longitude of the aircraft in decimal degrees
	Longitude float64

	// Registration number of the aircraft
	Registration string

	// Country of the aircraft based on the registration prefix
	Country string

	// Alitude of the aircraft in feet
	Altitude float64

	// Speed in knots
	Speed int

	// Distance between the specified location and the location of the aircraft in kilometers
	// Height is not taken into consideration
	Distance int

	// TrackerURL is to URL track the aircraft using the ADS-B website
	TrackerURL string

	// ImageThumbnailURL is the URL for the thumbnail of the aircraft
	ImageThumbnailURL string

	// ImageURL is the URL showing more images of the aircraft
	ImageURL string

	// Specifies the photographer of the image
	Photographer string

	// Percentage of cloud coverage, lower cloud coverage means that you have more chance to spot the aircraft
	CloudCoverage int

	// Bearing from your location to the aircraft
	BearingFromLocation float64

	// Bearing from the aircraft to your location
	BearingFromAircraft float64

	// Heading of the aircraft
	Heading float64

	// Specifies if it is a military type aircraft or not
	Military bool

	// Specifies if the aircraft is traveling towards your location or not
	Inbound bool

	// Specifies if the aircraft is on the ground
	OnGround bool

	// Airline of the aircraft
	Airline Airline

	// Origin of the flight
	Origin Airport

	// Destination of the flight
	Destination Airport
}

// FlightRouteResponse represents the structure of the response from the adsbdb.com API
type FlightRouteResponse struct {
	Response struct {
		FlightRoute FlightRoute `json:"flightroute"`
	} `json:"response"`
}

// FlightRoute contains detailed information about a flight route
type FlightRoute struct {
	Callsign     string  `json:"callsign"`
	CallsignICAO string  `json:"callsign_icao"`
	CallsignIATA string  `json:"callsign_iata"`
	Airline      Airline `json:"airline"`
	Origin       Airport `json:"origin"`
	Destination  Airport `json:"destination"`
}

// Airline contains information about an airline
type Airline struct {
	Name       string `json:"name"`
	ICAO       string `json:"icao"`
	IATA       string `json:"iata"`
	Country    string `json:"country"`
	CountryISO string `json:"country_iso"`
	Callsign   string `json:"callsign"`
}

// Airport contains information about an airport
type Airport struct {
	CountryISOName string  `json:"country_iso_name"`
	CountryName    string  `json:"country_name"`
	Elevation      int     `json:"elevation"`
	IATACode       string  `json:"iata_code"`
	ICAOCode       string  `json:"icao_code"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Municipality   string  `json:"municipality"`
	Name           string  `json:"name"`
}
