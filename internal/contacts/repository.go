package contacts

import (
	"database/sql"
)

// Repository defines methods for contact management
type IRepository interface {
	CreateContact(contact *Contact) (int, error)
	UpdateContact(contact *Contact) error
	SearchContacts(query string) ([]Contact, error)
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) IRepository {
	return &Repository{
		DB: db,
	}
}

// CreateContact stores a new contact with phone numbers
func (r *Repository) CreateContact(contact *Contact) (int, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Insert into contacts table
	var contactID int
	err = tx.QueryRow(`INSERT INTO contacts (first_name, last_name) VALUES ($1, $2) RETURNING id`,
		contact.FirstName, contact.LastName).Scan(&contactID)
	if err != nil {
		return 0, err
	}

	// Insert each phone number
	for _, number := range contact.PhoneNumbers {
		_, err = tx.Exec(`INSERT INTO phone_numbers (contact_id, number) VALUES ($1, $2)`, contactID, number)
		if err != nil {
			return 0, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return contactID, nil
}

// UpdateContact updates an existing contact and its phone numbers
func (r *Repository) UpdateContact(contact *Contact) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update contact details
	_, err = tx.Exec(`UPDATE contacts SET first_name = $1, last_name = $2 WHERE id = $3`,
		contact.FirstName, contact.LastName, contact.ID)
	if err != nil {
		return err
	}

	// Delete existing phone numbers
	_, err = tx.Exec(`DELETE FROM phone_numbers WHERE contact_id = $1`, contact.ID)
	if err != nil {
		return err
	}

	// Insert updated phone numbers
	for _, number := range contact.PhoneNumbers {
		_, err = tx.Exec(`INSERT INTO phone_numbers (contact_id, number) VALUES ($1, $2)`, contact.ID, number)
		if err != nil {
			return err
		}
	}

	// Commit transaction
	err = tx.Commit()
	return err
}

// SearchContacts finds contacts based on partial matches across all fields
func (r *Repository) SearchContacts(query string) ([]Contact, error) {
	rows, err := r.DB.Query(`
        SELECT c.id, c.first_name, c.last_name, p.number
        FROM contacts c
        LEFT JOIN phone_numbers p ON c.id = p.contact_id
        WHERE c.first_name ILIKE '%' || $1 || '%'
           OR c.last_name ILIKE '%' || $1 || '%'
           OR p.number ILIKE '%' || $1 || '%'
    `, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map contacts by ID for easier aggregation of phone numbers
	contactsMap := make(map[int]*Contact)
	for rows.Next() {
		var contactID int
		var firstName, lastName, phoneNumber string
		err = rows.Scan(&contactID, &firstName, &lastName, &phoneNumber)
		if err != nil {
			return nil, err
		}

		contact, exists := contactsMap[contactID]
		if !exists {
			contact = &Contact{
				ID:           contactID,
				FirstName:    firstName,
				LastName:     lastName,
				PhoneNumbers: []string{},
			}
			contactsMap[contactID] = contact
		}
		contact.PhoneNumbers = append(contact.PhoneNumbers, phoneNumber)
	}

	// Convert map to slice
	contacts := make([]Contact, 0, len(contactsMap))
	for _, contact := range contactsMap {
		contacts = append(contacts, *contact)
	}
	return contacts, nil
}
