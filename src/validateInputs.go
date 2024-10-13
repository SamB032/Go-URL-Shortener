package main

import "regexp"

// Checks if the input string is in the form of a URL
func validateIfURL(url string) bool {
  // Regular expression to match URLs with or without schemes and with sub-directories
	re := regexp.MustCompile(`^(https?://)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(/[\w\-._~:/?#[\]@!$&'()*+,;=]*)?$`)

  // Return true if the URL matches the regex, false otherwise
	return re.MatchString(url)
}
