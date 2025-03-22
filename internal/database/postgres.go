package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PGClient wraps the sqlx.DB instance for better resource management.
type PGClient struct {
	DB *sqlx.DB
}

// Connect establishes a connection to the database and returns a PGClient instance.
func Connect(provider, user, password, dbname, host string, sslmode bool) (*PGClient, error) {
	sslmodeStr := "disable"
	if sslmode {
		sslmodeStr = "require"
	}

	dsn := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=%s", provider, user, password, host, dbname, sslmodeStr)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	log.Println("Successfully connected to the database")
	return &PGClient{DB: db}, nil
}

// Close closes the database connection.
func (dbw *PGClient) Close() {
	if err := dbw.DB.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	} else {
		log.Println("Database connection closed")
	}
}
