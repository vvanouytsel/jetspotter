package aircraft

import (
	"context"
	"database/sql"
	"fmt"
)

type Aircraft struct {
	ID          string
	Callsign    string
	Type        string
	Description string
	TailNumber  string
	ICAO        string
	ImageURL    string
	IsMilitary  bool
}
type AircraftRepo interface {
	GetByID(aircraftID string) (*Aircraft, error)
}

type PostgresAircaftRepo struct {
	db *sql.DB
}

func NewPostgresAircraftRepo(db *sql.DB) *PostgresAircaftRepo {
	return &PostgresAircaftRepo{
		db: db,
	}
}

func (pr *PostgresAircaftRepo) GetByID(aircraftID string) (*Aircraft, error) {
	// Logic to get a by ID goes here.
	var a Aircraft
	row := pr.db.QueryRowContext(context.TODO(), "SELECT * FROM aircraft WHERE aircraft_id = $1", aircraftID)
	err := row.Scan(&a.ID, &a.Callsign, &a.Type, &a.Description, &a.TailNumber, &a.ICAO, &a.ImageURL, &a.IsMilitary)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

type AircraftService struct {
	repo AircraftRepo
}

func NewAircraftService(ar AircraftRepo) *AircraftService {
	return &AircraftService{
		repo: ar,
	}
}

func (as AircraftService) VerifyAircraftMilitary(ac Aircraft) error {
	// Business logic
	aircraft, err := as.repo.GetByID(ac.ID)
	if err != nil {
		return err
	}

	if aircraft == nil {
		fmt.Println("No aircraft found.")
		return nil
	}

	fmt.Printf("Callsign: %s", aircraft.Callsign)
	if aircraft.IsMilitary {
		fmt.Println("Aircraft is military")
		return nil
	}

	return nil
}
