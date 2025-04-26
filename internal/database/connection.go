package database

import (
	_ "github.com/lib/pq"

	"database/sql"
	"fmt"
	"log/slog"
)

type Connection struct {
	connection *sql.DB
}

type DBInterface interface {
	CheckShortkeyExists(shortKey string) (bool, error)
	CheckIfURLExists(url string) (bool, error)
	AddRecord(oldurl string, shortKey string) error
	FindURLUsingShortkey(shortKey string) (string, error)
	FindShortkeyUsingURL(url string) (string, error)
}

func ConnectToDatabase(pgHost, pgPort, pgUser, pgPassword, pgName string, logger *slog.Logger) *Connection {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPassword, pgName)

	//Initalise a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Error("Failed to connect to database",
			slog.String("host", pgHost),
			slog.String("DBName", pgName),
			slog.String("Error", err.Error()),
		)
		return nil
	}

	// Ensure the connection is properly closed on error
	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	// Make a ping to the database to see if its alive
	err = db.Ping()
	if err != nil {
		logger.Error("Failed to connect to ping database",
			slog.String("host", pgHost),
			slog.String("DBName", pgName),
			slog.String("Error", err.Error()),
		)
		return nil
	}

	logger.Error("Successfully connect to the database",
		slog.String("host", pgHost),
		slog.String("DBName", pgName),
	)

	return &Connection{connection: db}
}
