package api

// album represents data about a record album.
type Aircraft struct {
	ID          string `json:"id"`
	Callsign    string `json:"callsign"`
	Type        string `json:"type"`
	Description string `json:"description"`
	TailNumber  string `json:"tailNumber"`
	ICAO        string `json:"ICAO"`
	ImageURL    string `json:"imageURL"`
	IsMilitary  bool   `json:"isMilitary"`
}
