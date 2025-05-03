package database

import (
	_ "github.com/lib/pq"

	"database/sql"
	"fmt"
	"log/slog"
)

var SqlOpen = sql.Open

type DBInterface interface {
	CheckShortkeyExists(shortKey string) (bool, error)
	CheckIfURLExists(url string) (bool, error)
	AddRecord(oldurl string, shortKey string) error
	FindURLUsingShortkey(shortKey string) (string, error)
	FindShortkeyUsingURL(url string) (string, error)
}

type Connection struct {
	Connection *sql.DB
}

func NewConnection(db *sql.DB, logger *slog.Logger, host, dbname string) *Connection {
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database",
			slog.String("host", host),
			slog.String("DBName", dbname),
			slog.String("Error", err.Error()),
		)
		db.Close()
		return nil
	}

	logger.Info("Successfully connected to the database",
		slog.String("host", host),
		slog.String("DBName", dbname),
	)

	return &Connection{Connection: db}
}

func ConnectToDatabase(pgHost, pgPort, pgUser, pgPassword, pgName string, logger *slog.Logger) *Connection {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPassword, pgName)

	db, err := SqlOpen("postgres", psqlInfo)
	if err != nil {
		logger.Error("Failed to open connection",
			slog.String("Error", err.Error()),
		)
		return nil
	}

	return NewConnection(db, logger, pgHost, pgName)
}
