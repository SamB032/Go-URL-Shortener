package main

import (
  _ "github.com/lib/pq"

  "log"
  "time"
)

func checkShortkeyExists(shortKey string) bool {
  query := `SELECT EXISTS(SELECT 1 FROM url WHERE shortkey = $1)`

  var exists bool
  err := dbConnection.QueryRow(query, shortKey).Scan(&exists)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
  return exists
}

func checkIfURLExists(url string) bool {
  query := `SELECT EXISTS(SELECT 1 from URL where old_url = $1)`

  var exists bool
  err := dbConnection.QueryRow(query, url).Scan(&exists)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
  return exists
}

func addRecord(oldurl string, shortKey string) error {
  timestamp := time.Now() //Record timestamp of when record is added

  _, err := dbConnection.Exec(`INSERT INTO url (created_at, old_url, shortkey) VALUES ($1, $2, $3)`, timestamp, oldurl, shortKey)

  if err != nil {
    return err
	}
  return nil
}

// Find the corresponding short url when given a shortKey
func findURLUsingShortkey(shortKey string) (string, error) {
  var oldURL string

  query := "SELECT old_url FROM url WHERE shortkey = $1"
  err := dbConnection.QueryRow(query, shortKey).Scan(&oldURL)
  if err != nil {
      return "", err // Return any other errors
  }
  return oldURL, nil
}
