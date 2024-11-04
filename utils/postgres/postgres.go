package postgres

import (
	"database/sql"
	"fmt"
	"io/fs"

	"phonebook/internal/migrations"
	"phonebook/utils/configs"

	_ "github.com/jackc/pgx/v5/stdlib" // Postgres driver
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
)

// PostgresConfig holds PostgreSQL connection parameters
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// DefaultPostgresConfig loads default PostgreSQL configuration
func DefaultPostgresConfig() PostgresConfig {
	dbConfig := configs.C().PSQL

	return PostgresConfig{
		Host:     dbConfig.Host,
		Port:     dbConfig.Port,
		User:     dbConfig.User,
		Password: dbConfig.Password,
		Database: dbConfig.Database,
		SSLMode:  dbConfig.SSLMode,
	}
}

// String returns the DSN (Data Source Name) for connecting to PostgreSQL
func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)
}

// Postgres struct holds the PostgreSQL database connection
type Postgres struct {
	DB *sql.DB
}

// Global Postgres instance
var PostgresInstance Postgres

// ConnectPostgres initializes a PostgreSQL connection using the provided config
func (p *Postgres) ConnectPostgres(cfg PostgresConfig) error {
	var err error
	p.DB, err = sql.Open("pgx", cfg.String())
	if err != nil {
		log.Error().Err(err).Msg("Failed to open PostgreSQL connection")
		return err
	}

	// Verify the connection
	if err = p.DB.Ping(); err != nil {
		log.Error().Err(err).Msg("Failed to ping PostgreSQL")
		return err
	}

	log.Info().Msg("Connected to PostgreSQL!")
	return nil
}

// DisconnectPostgres gracefully disconnects the PostgreSQL client
func (p *Postgres) DisconnectPostgres() {
	if err := p.DB.Close(); err != nil {
		log.Error().Err(err).Msg("Error disconnecting from PostgreSQL")
	} else {
		log.Info().Msg("Disconnected from PostgreSQL!")
	}
}

// Migrate applies database migrations from the given directory
func Migrate(db *sql.DB, dir string) {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal().Err(err).Msg("Failed to set PostgreSQL dialect")
	}

	if err := goose.Up(db, dir); err != nil {
		log.Fatal().Err(err).Msg("Migration failed")
	}
}

// MigrateFS applies database migrations using an in-memory filesystem (for embedded migrations)
func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) {
	if dir == "" {
		dir = "."
	}
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	Migrate(db, dir)
}

// RunPostgres initializes the PostgreSQL connection and applies migrations
// Don't forget to call this function in the main function and defer the DisconnectPostgres function
func RunPostgres() {
	// Load the default PostgreSQL configuration
	config := DefaultPostgresConfig()

	// Initialize PostgreSQL using the config
	if err := PostgresInstance.ConnectPostgres(config); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}

	// Apply migrations
	MigrateFS(PostgresInstance.DB, migrations.FS, ".")
}
