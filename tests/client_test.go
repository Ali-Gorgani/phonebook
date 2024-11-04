package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"phonebook/internal/api-gateway/http/client"
	"phonebook/internal/contacts"
	"testing"
)

// Helper function to create a mock server for each test case
func setupMockServer(t *testing.T, handlerFunc http.HandlerFunc) (*httptest.Server, *client.Client) {
	ts := httptest.NewServer(handlerFunc)
	t.Cleanup(func() { ts.Close() }) // Ensure server closes after test
	return ts, client.NewClient(ts.URL)
}

func TestClientAddContact(t *testing.T) {
	_, c := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/contacts" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		var contact contacts.Contact
		if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Respond with a mock ID
		response := map[string]int{"contact_id": 1}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	})

	// Prepare test contact and make request
	contact := contacts.Contact{FirstName: "John", LastName: "Doe", PhoneNumbers: []string{"1234567890"}}
	id, err := c.AddContact(&contact)
	if err != nil {
		t.Fatalf("AddContact failed: expected no error, got %v", err)
	}

	if id != 1 {
		t.Fatalf("AddContact failed: expected contact ID to be 1, got %d", id)
	}
}

func TestClientUpdateContact(t *testing.T) {
	_, c := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/contacts/1" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		var contact contacts.Contact
		if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Respond with success status
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
	})

	// Prepare test contact and make request
	contact := contacts.Contact{ID: 1, FirstName: "Jane", LastName: "Doe", PhoneNumbers: []string{"0987654321"}}
	err := c.UpdateContact(&contact)
	if err != nil {
		t.Fatalf("UpdateContact failed: expected no error, got %v", err)
	}
}

func TestClientSearchContacts(t *testing.T) {
	_, c := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/contacts/search" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		query := r.URL.Query().Get("q") // Fix: Get the correct query parameter
		if query == "" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Mock response with a JSON structure wrapping contacts
		response := struct {
			Contacts []contacts.Contact `json:"contacts"`
		}{
			Contacts: []contacts.Contact{
				{ID: 1, FirstName: "John", LastName: "Doe", PhoneNumbers: []string{"1234567890"}},
				{ID: 2, FirstName: "Jane", LastName: "Doe", PhoneNumbers: []string{"0987654321"}},
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response) // Encode the response as a JSON object
	})

	// Make request to search contacts
	contacts, err := c.SearchContacts("Doe")
	if err != nil {
		t.Fatalf("SearchContacts failed: expected no error, got %v", err)
	}

	// Validate response
	if len(contacts) != 2 {
		t.Fatalf("SearchContacts failed: expected 2 contacts, got %d", len(contacts))
	}

	if contacts[0].FirstName != "John" || contacts[1].FirstName != "Jane" {
		t.Fatalf("SearchContacts failed: unexpected contacts returned")
	}
}
