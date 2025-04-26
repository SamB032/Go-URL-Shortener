package database

import (
	_ "github.com/lib/pq"

	"fmt"
	"time"
)

// Checks whether a shortkey already exists in the database
func (db *Connection) CheckShortkeyExists(shortKey string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM url WHERE shortkey = $1)`

	var exists bool
	err := db.connection.QueryRow(query, shortKey).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error executing query: %v", err)
	}
	return exists, nil
}

// Check whether a URL already exists in the database
func (db *Connection) CheckIfURLExists(url string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 from URL where old_url = $1)`

	var exists bool
	err := db.connection.QueryRow(query, url).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// Adds a mapping shortkey <-> url to the database
func (db *Connection) AddRecord(oldurl string, shortKey string) error {
	timestamp := time.Now() //Record timestamp of when record is added

	_, err := db.connection.Exec(`INSERT INTO url (created_at, old_url, shortkey) VALUES ($1, $2, $3)`, timestamp, oldurl, shortKey)

	if err != nil {
		return err
	}
	return nil
}

// Find the corresponding short url when given a shortKey
func (db *Connection) FindURLUsingShortkey(shortKey string) (string, error) {
	var oldURL string

	query := "SELECT old_url FROM url WHERE shortkey = $1"
	err := db.connection.QueryRow(query, shortKey).Scan(&oldURL)
	if err != nil {
		return "", err
	}
	return oldURL, nil
}

// Find shortkey when given the url
func (db *Connection) FindShortkeyUsingURL(url string) (string, error) {
	var shortKey string

	query := "SELECT shortkey FROM url WHERE old_url = $1"
	err := db.connection.QueryRow(query, url).Scan(&shortKey)
	if err != nil {
		return "", err
	}
	return shortKey, nil
}
