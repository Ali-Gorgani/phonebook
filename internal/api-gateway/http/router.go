package http

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()
	handler := NewHandler(db)

	// Define routes
	router.POST("/contacts", handler.CreateContactHandler)
	router.PUT("/contacts/:id", handler.UpdateContactHandler)
	router.GET("/contacts/search", handler.SearchContactsHandler)

	return router
}
