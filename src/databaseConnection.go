package main

import (
	_ "github.com/lib/pq"

	"database/sql"
  "log/slog"
	"fmt"
	"os"
)

type PostgresData struct {
  host     string
  port     string
  user     string
  password string
  dbName   string
}

type DBConnection struct {
  connection *sql.DB
}

//Load the posgres data from environment variables and return it as a struct
func getDatabaseInfo() PostgresData {
  return PostgresData{
    os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"),
  }
}

func connectToDatabase(logger *slog.Logger) (*DBConnection, error) {
  dbInfo := getDatabaseInfo()

  psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
    dbInfo.host, dbInfo.port, dbInfo.user, dbInfo.password, dbInfo.dbName)

  //Initalise a connection to the database
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    logger.Error("Failed to connect to database",
      slog.String("host", dbInfo.host),
      slog.String("DBName", dbInfo.dbName),
      slog.String("Error", err.Error()),
    )
    return nil, fmt.Errorf("database connection failed: %w", err)
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
    logger.Info("Failed to connect to ping database",
      slog.String("host", dbInfo.host),
      slog.String("DBName", dbInfo.dbName),
      slog.String("Error", err.Error()),
    )
    return nil, fmt.Errorf("database ping failed: %w", err)
  }

  logger.Info("Successfully connect to the database",
    slog.String("host", dbInfo.host),
    slog.String("DBName", dbInfo.dbName),
  )

  //Save the connector as a struct
  return &DBConnection{db}, nil
}
