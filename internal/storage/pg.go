package storage

import (
	"TestWebServer/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Importing the PostgreSQL driver.
)

// Database struct encapsulates the database connection.
type Database struct {
	pg *sqlx.DB // Pointer to the sqlx.DB instance for database operations.
}

// InitDatabase initializes a new database connection using the provided configuration.
// It returns a pointer to the sqlx.DB instance and an error if any occurs during the connection process.
func InitDatabase(cfg config.Config) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", cfg.DB) // Connecting to the PostgreSQL database.
	if err != nil {
		return nil, err
	}
	err = conn.Ping() // Pinging the database to ensure the connection is alive.
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// NewDatabase creates a new Database instance with the provided sqlx.DB connection.
func NewDatabase(conn *sqlx.DB) *Database {
	return &Database{pg: conn}
}
