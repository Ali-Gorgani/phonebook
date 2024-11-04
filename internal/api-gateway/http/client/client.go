package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"phonebook/internal/contacts"
)

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

// AddContact sends a request to create a new contact
func (c *Client) AddContact(contact *contacts.Contact) (int, error) {
	data, err := json.Marshal(contact)
	if err != nil {
		return 0, err
	}

	resp, err := http.Post(c.baseURL+"/contacts", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return 0, errors.New("failed to create contact")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	id, ok := result["contact_id"].(float64)
	if !ok {
		return 0, errors.New("invalid response format")
	}

	return int(id), nil
}

// UpdateContact sends a request to update an existing contact
func (c *Client) UpdateContact(contact *contacts.Contact) error {
	data, err := json.Marshal(contact)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/contacts/%d", c.baseURL, contact.ID), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to update contact")
	}

	return nil
}

// SearchContacts sends a request to search for contacts by query
func (c *Client) SearchContacts(query string) ([]contacts.Contact, error) {
	resp, err := http.Get(fmt.Sprintf("%s/contacts/search?q=%s", c.baseURL, query))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("search failed")
	}

	var result struct {
		Contacts []contacts.Contact `json:"contacts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Contacts, nil
}
