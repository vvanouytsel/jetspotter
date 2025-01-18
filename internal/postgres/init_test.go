package postgres

import (
	"database/sql"
	"fmt"
	"testing"
)

const (
	totalTables = 7
)

func TestInit(t *testing.T) {
	err := initialize()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	host := "localhost"
	port := 5432
	user := "jetspotter"
	password := "dev"
	dbname := "jetspotter"

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer db.Close()
	tables := 0
	if err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';").Scan(&tables); err != nil {
		t.Errorf("Error: %v", err)
	}

	if tables != totalTables {
		t.Errorf("Expected %v tables but found %v", totalTables, tables)
	}
}
