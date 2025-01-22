package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var (
	initQueries = []string{`
		CREATE TABLE IF NOT EXISTS aircraft (
			aircraft_id SERIAL PRIMARY KEY,
			callsign VARCHAR(50),
			type VARCHAR(100),
			description VARCHAR(255),
			tail_number VARCHAR(50),
			icao VARCHAR(50),
			image_url VARCHAR(255),
			military BOOLEAN
		);`,
		`CREATE TABLE IF NOT EXISTS spot_configurations (
			spot_configuration_id SERIAL PRIMARY KEY,
			lattitude NUMERIC,
			longitude NUMERIC,
			max_range_kilometers INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS spots (
			spot_id SERIAL PRIMARY KEY,
			aircraft_id int REFERENCES aircraft(aircraft_id),
			spot_configuration_id int REFERENCES spot_configurations(spot_configuration_id),
			spot_date DATE
		);`,
		`CREATE TABLE IF NOT EXISTS notification_types (
			notification_type_id SERIAL PRIMARY KEY,
			name VARCHAR(50)
		);`,
		`CREATE TABLE IF NOT EXISTS notifications (
			notification_id SERIAL PRIMARY KEY,
			spot_id int REFERENCES spots(spot_id),
			notification_type_id int REFERENCES notification_types(notification_type_id)
		);`,
		`CREATE TABLE IF NOT EXISTS notification_configurations (
			notification_configuration_id SERIAL PRIMARY KEY,
			notification_type_id int REFERENCES notification_types(notification_type_id),
			webhook_url VARCHAR(255),
			gotify_token VARCHAR(255),
			ntfy_topic VARCHAR(255),
			ntfy_server VARCHAR(255),
			max_range_kilometers INTEGER,
			max_altitude_feet INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL PRIMARY KEY,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			email VARCHAR(100),
			password VARCHAR(255),
			notification_configuration_id int REFERENCES notification_configurations(notification_configuration_id)
		);`,
		`CREATE EXTENSION IF NOT EXISTS pgcrypto;`,
	}
)

// Initialize creates the tables in the database if they did not exist yet.
func Initialize() error {
	host := "localhost"
	port := 5432
	user := "jetspotter"
	password := "dev"
	dbname := "jetspotter"

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Connected to the database!")

	for _, query := range initQueries {
		_, err := db.Query(query)
		if err != nil {
			fmt.Println("Failed query: ", query)
			return err
		}

	}

	return nil
}
