package http

import (
	"database/sql"
	"log"
	"net/http"
	"phonebook/internal/contacts"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *contacts.Service
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		service: contacts.NewService(db),
	}
}

// CreateContactHandler handles the creation of a new contact
func (h *Handler) CreateContactHandler(c *gin.Context) {
	var contact contacts.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	contactID, err := h.service.CreateContact(&contact)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create contact"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"contact_id": contactID})
}

// UpdateContactHandler handles updating an existing contact
func (h *Handler) UpdateContactHandler(c *gin.Context) {
	var contact contacts.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Get contact ID from URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID"})
		return
	}
	contact.ID = id

	if err := h.service.UpdateContact(&contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact updated successfully"})
}

// SearchContactsHandler handles searching for contacts
func (h *Handler) SearchContactsHandler(c *gin.Context) {
	query := c.Query("q")
	contacts, err := h.service.SearchContacts(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not search contacts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"contacts": contacts})
}
