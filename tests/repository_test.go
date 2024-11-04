package tests

import (
	"phonebook/internal/contacts"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateContactWithTransaction(t *testing.T) {
	repo := contacts.NewRepository(mockDB)

	contact := &contacts.Contact{FirstName: "John", LastName: "Doe", PhoneNumbers: []string{"1234567890"}}

	// Set expectations for the mocked transaction
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO contacts \(first_name, last_name\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs(contact.FirstName, contact.LastName).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) // Return ID 1
	mock.ExpectExec(`INSERT INTO phone_numbers \(contact_id, number\) VALUES \(\$1, \$2\)`).
		WithArgs(1, "1234567890").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test CreateContact method
	id, err := repo.CreateContact(contact)
	if err != nil {
		t.Fatalf("CreateContact failed, expected no error, got %v", err)
	}

	if id != 1 {
		t.Fatalf("expected contact ID to be 1, got %d", id)
	}

	// Verify all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestUpdateContactWithTransaction(t *testing.T) {
	repo := contacts.NewRepository(mockDB)

	contact := &contacts.Contact{ID: 1, FirstName: "Jane", LastName: "Doe", PhoneNumbers: []string{"0987654321"}}

	// Set expectations for the mocked transaction
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE contacts SET first_name = \$1, last_name = \$2 WHERE id = \$3`).
		WithArgs(contact.FirstName, contact.LastName, contact.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`DELETE FROM phone_numbers WHERE contact_id = \$1`).
		WithArgs(contact.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(`INSERT INTO phone_numbers \(contact_id, number\) VALUES \(\$1, \$2\)`).
		WithArgs(contact.ID, "0987654321").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Test UpdateContact method
	err := repo.UpdateContact(contact)
	if err != nil {
		t.Fatalf("UpdateContact failed, expected no error, got %v", err)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestSearchContactsWithTransaction(t *testing.T) {
	repo := contacts.NewRepository(mockDB)

	// Define expected mock behavior and rows to return
	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "number"}).
		AddRow(1, "John", "Doe", "1234567890").
		AddRow(2, "Jane", "Doe", "0987654321")

	// Set expectations for the mocked transaction
	mock.ExpectQuery(`SELECT c\.id, c\.first_name, c\.last_name, p\.number FROM contacts c LEFT JOIN phone_numbers p ON c\.id = p\.contact_id WHERE c\.first_name ILIKE '%' \|\| \$1 \|\| '%' OR c\.last_name ILIKE '%' \|\| \$1 \|\| '%' OR p\.number ILIKE '%' \|\| \$1 \|\| '%'`).
		WithArgs("Doe").
		WillReturnRows(rows)

	// Test SearchContacts method
	contacts, err := repo.SearchContacts("Doe")
	if err != nil {
		t.Fatalf("SearchContacts failed, expected no error, got %v", err)
	}

	// Validate results
	if len(contacts) != 2 {
		t.Fatalf("expected 2 contacts, got %d", len(contacts))
	}
	if contacts[0].FirstName != "John" || contacts[1].FirstName != "Jane" {
		t.Fatalf("unexpected contacts returned, got %+v", contacts)
	}

	// Verify expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
