package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// Declare global variables for the mock database
var (
	mockDB *sql.DB
	mock   sqlmock.Sqlmock
)

func TestMain(m *testing.M) {
	// Initialize the mock database
	var err error
	mockDB, mock, err = sqlmock.New()
	if err != nil {
		log.Fatalf("failed to create mock db: %v", err)
	}
	defer mockDB.Close() // Ensure mockDB is closed after tests

	// Run the tests
	code := m.Run()

	log.Println("Tests completed.")
	os.Exit(code)
}
