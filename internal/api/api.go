package api

import (
	"database/sql"
	"fmt"
	"jetspotter/internal/aircraft"
	"jetspotter/internal/configuration"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var ac = []aircraft.Aircraft{
	{
		ID:          "1",
		Callsign:    "Apex1",
		Type:        "Cessna",
		Description: "A small, single-engine aircraft that is perfect for short trips.",
		TailNumber:  "N12345",
		ICAO:        "C172",
		ImageURL:    "https://example.com/image.jpg",
		IsMilitary:  false,
	},
	{
		ID:          "2",
		Callsign:    "Viper1",
		Type:        "F-16",
		Description: "A multirole fighter jet that is used by many countries around the world.",
		TailNumber:  "AF12345",
		ICAO:        "F16",
		ImageURL:    "https://example.com/image.jpg",
		IsMilitary:  true,
	},
}

func getAllAircraft(c *gin.Context) {
	fmt.Println("Getting all aircraft")
}

func getAircraftByID(c *gin.Context) {
	host := "localhost"
	port := 5432
	user := "jetspotter"
	password := "dev"
	dbname := "jetspotter"

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Get aircraft by id")
	postgresAircraftRepo := aircraft.NewPostgresAircraftRepo(db)
	aircraftService := aircraft.NewAircraftService(postgresAircraftRepo)
	err = aircraftService.VerifyAircraftMilitary(ac[0])
	if err != nil {
		log.Fatalf("Error verifying aircraft: %v\n", err)
	}

}

func HandleAPI(config configuration.Config) error {
	router := gin.Default()
	router.GET("/aircraft", getAllAircraft)
	router.GET("/aircraft/:id", getAircraftByID)
	err := router.Run("localhost:8080")
	if err != nil {
		return err
	}
	return nil
}
