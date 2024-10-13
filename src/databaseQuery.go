package main

import (
  _ "github.com/lib/pq"

  "log"
)

func checkShortkeyExists(shortKey string) bool {
  query := `SELECT EXISTS(SELECT 1 FROM url WHERE new_url = $1)`

  var exists bool
  err := dbConnection.QueryRow(query, shortKey).Scan(&exists)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}

  return exists
}
