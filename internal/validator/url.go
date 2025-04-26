package validate

import (
	"regexp"

	database "github.com/SamB032/Go-URL-Shortener/internal/database"
)

const urlRegex = `(?i)(https?://(?:www\.)?[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]\.[^\s]{2,}|https?://(?:www\.)?[a-zA-Z0-9]+\.[^\s]{2,}|www\.[a-zA-Z0-9]+\.[^\s]{2,})`

// Checks if the input string is in the form of a URL, this also performs a database call to check if url already exists
func ValidateURL(url string, db database.DBInterface) (bool, bool, error) {
	exists, err := db.CheckIfURLExists(url)
	if err != nil {
		return false, false, err
	}
	if exists {
		return true, true, nil
	}

	// Regular expression to match URLs with or without schemes and with sub-directories
	re := regexp.MustCompile(urlRegex)

	// Return true if the url is valid
	return len(url) <= 255 && re.MatchString(url), false, nil
}
