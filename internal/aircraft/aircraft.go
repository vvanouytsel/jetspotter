package aircraft

// Type of aircraft
type Type struct {
	// Identifier of the type
	Identifier string
	// Description of the type
	Description string
}

var (
	// B77L BOEING 777-200LR
	B77L = Type{
		Identifier:  "B77L",
		Description: "BOEING 777-200LR",
	}

	// PC12 PILATUS PC-12
	PC12 = Type{
		Identifier:  "PC12",
		Description: "PILATUS PC-12",
	}

	// B789 BOEING 787-9 Dreamliner
	B789 = Type{
		Identifier:  "B789",
		Description: "BOEING 787-9 Dreamliner",
	}

	// F16 GENERAL DYNAMICS F-16 Fighting Falcon
	F16 = Type{
		Identifier:  "F16",
		Description: "GENERAL DYNAMICS F-16 Fighting Falcon",
	}

	// E295 EMBRAER ERJ-190-400
	E295 = Type{
		Identifier:  "E295",
		Description: "EMBRAER ERJ-190-400",
	}

	// A320 AIRBUS A-320
	A320 = Type{
		Identifier:  "A320",
		Description: "AIRBUS A-320",
	}

	// C550 CESSNA 550 Citation S2
	C550 = Type{
		Identifier:  "C550",
		Description: "CESSNA 550 Citation S2",
	}

	// E170 EMBRAER ERJ-170-100
	E170 = Type{
		Identifier:  "E170",
		Description: "EMBRAER ERJ-170-100",
	}

	// A321 AIRBUS A-321
	A321 = Type{
		Identifier:  "A321",
		Description: "AIRBUS A-321",
	}

	// A319 AIRBUS A-319
	A319 = Type{
		Identifier:  "A319",
		Description: "AIRBUS A-319",
	}

	// A400 ATLAS - A400
	A400 = Type{
		Identifier:  "A400",
		Description: "A400 ATLAS",
	}

	// ALL ALL
	ALL = Type{
		Identifier:  "ALL",
		Description: "ALL",
	}

	// MILITARY MILITIARY
	MILITARY = Type{
		Identifier:  "MILITARY",
		Description: "MILITARY",
	}
)
