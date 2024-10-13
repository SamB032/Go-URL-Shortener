package main

import "regexp"

// Checks if the input string is in the form of a URL, this also performs a database call to check if url already exists
func validateIfURL(url string) bool {
  // Regular expression to match URLs with or without schemes and with sub-directories
  re := regexp.MustCompile(`(?i)(https?://(?:www\.)?[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]\.[^\s]{2,}|www\.[a-zA-Z0-9][a-zA-Z0-9-]*[a-zA-Z0-9]\.[^\s]{2,}|https?://(?:www\.)?[a-zA-Z0-9]+\.[^\s]{2,}|www\.[a-zA-Z0-9]+\.[^\s]{2,})`)

  // Return true if the url is valid
  return len(url) <= 255 && re.MatchString(url) && !checkIfURLExists(url)
}
