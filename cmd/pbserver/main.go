package main

import (
	"os"
	"phonebook/internal/api-gateway/http"
	"phonebook/utils/configs"
	"phonebook/utils/postgres"

	"github.com/rs/zerolog"
)

func main() {
	configs.RunConfig(".")
	postgres.RunPostgres()

	// Initialize router
	db := postgres.PostgresInstance.DB
	router := http.NewRouter(db)

	// Start the server on port 1234 using zerolog
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("Starting server on port 1234...")
	if err := router.Run(":1234"); err != nil {
		logger.Fatal().Err(err).Msg("Could not start server")
	}
}
