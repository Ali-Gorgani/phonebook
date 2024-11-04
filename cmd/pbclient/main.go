package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"phonebook/internal/api-gateway/http/client"
	"phonebook/internal/contacts"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: pbclient <baseURL>")
	}

	baseURL := os.Args[1]
	c := client.NewClient(baseURL)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Phone Book Client. Type 'help' for commands.")

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		args := strings.Split(input, " ")
		command := args[0]

		switch command {
		case "add":
			addContact(c, args[1:])
		case "update":
			updateContact(c, args[1:])
		case "search":
			searchContacts(c, args[1:])
		case "help":
			printHelp()
		case "exit":
			fmt.Println("Exiting client.")
			return
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}
}

func addContact(c *client.Client, args []string) {
	if len(args) < 3 {
		fmt.Println("Usage: add <first_name> <last_name> <phone_numbers...>")
		return
	}

	firstName, lastName := args[0], args[1]
	phoneNumbers := args[2:]

	contact := contacts.Contact{
		FirstName:    firstName,
		LastName:     lastName,
		PhoneNumbers: phoneNumbers,
	}

	id, err := c.AddContact(&contact)
	if err != nil {
		fmt.Printf("Error adding contact: %v\n", err)
		return
	}

	fmt.Printf("Contact created with ID %d\n", id)
}

func updateContact(c *client.Client, args []string) {
	if len(args) < 4 {
		fmt.Println("Usage: update <id> <first_name> <last_name> <phone_numbers...>")
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	firstName, lastName := args[1], args[2]
	phoneNumbers := args[3:]

	contact := contacts.Contact{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		PhoneNumbers: phoneNumbers,
	}

	if err := c.UpdateContact(&contact); err != nil {
		fmt.Printf("Error updating contact: %v\n", err)
		return
	}

	fmt.Println("Contact updated successfully")
}

func searchContacts(c *client.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: search <query>")
		return
	}

	query := args[0]
	contacts, err := c.SearchContacts(query)
	if err != nil {
		fmt.Printf("Error searching contacts: %v\n", err)
		return
	}

	fmt.Println("Search Results:")
	for _, contact := range contacts {
		fmt.Printf("ID: %d, Name: %s %s, Phones: %v\n", contact.ID, contact.FirstName, contact.LastName, contact.PhoneNumbers)
	}
}

func printHelp() {
	fmt.Println("Commands:")
	fmt.Println("  add <first_name> <last_name> <phone_numbers...> - Add a new contact")
	fmt.Println("  update <id> <first_name> <last_name> <phone_numbers...> - Update a contact")
	fmt.Println("  search <query> - Search contacts by name or phone number")
	fmt.Println("  help - Show available commands")
	fmt.Println("  exit - Exit the client")
}
