package aircraft

type AircraftType struct {
	Identifier  string
	Description string
}

var (
	B77L = AircraftType{
		Identifier:  "B77L",
		Description: "BOEING 777-200LR",
	}

	PC12 = AircraftType{
		Identifier:  "PC12",
		Description: "PILATUS PC-12",
	}

	B789 = AircraftType{
		Identifier:  "B789",
		Description: "BOEING 787-9 Dreamliner",
	}

	F16 = AircraftType{
		Identifier:  "F16",
		Description: "GENERAL DYNAMICS F-16 Fighting Falcon",
	}

	E295 = AircraftType{
		Identifier:  "E295",
		Description: "EMBRAER ERJ-190-400",
	}

	A320 = AircraftType{
		Identifier:  "A320",
		Description: "AIRBUS A-320",
	}

	C550 = AircraftType{
		Identifier:  "C550",
		Description: "CESSNA 550 Citation S2",
	}

	E170 = AircraftType{
		Identifier:  "E170",
		Description: "EMBRAER ERJ-170-100",
	}

	A321 = AircraftType{
		Identifier:  "A321",
		Description: "AIRBUS A-321",
	}

	A319 = AircraftType{
		Identifier:  "A319",
		Description: "AIRBUS A-319",
	}

	ALL = AircraftType{
		Identifier:  "ALL",
		Description: "ALL",
	}
)
