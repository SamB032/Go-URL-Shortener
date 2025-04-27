package shortkey

import (
	"math/rand"

	database "github.com/SamB032/Go-URL-Shortener/internal/database"
)

const URL_CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const URL_LENGTH = 5

// Generate a shortkey, return only if one is found that does not already exists
func CreateShortKey(dbConnection database.DBInterface) (string, error) {
	var newurl string
	for {
		newurl = GenerateShortKey()
		exists, err := dbConnection.CheckShortkeyExists(newurl)
		if err != nil {
			return "", err
		} else if !exists {
			return newurl, nil
		}
	}
}

// Generates a random key of URL_LENGTH and contains URL_CHARSET
func GenerateShortKey() string {
	shortKey := make([]byte, URL_LENGTH)

	for i := range shortKey {
		shortKey[i] = URL_CHARSET[rand.Intn(len(URL_CHARSET))]
	}

	return string(shortKey)
}
