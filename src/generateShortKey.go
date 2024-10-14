package main

import "math/rand"

const URL_CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const URL_LENGTH = 5

//Generate a shortkey, return only if one is found that does not already exists
func createShortKey() (string, error) {
  var newurl string
  for {
    newurl = generateShortKey()
    exists, err := dbConnection.checkShortkeyExists(newurl)
    if err != nil {
      return "", err
    } else if !exists {
      return newurl, nil
    }
  } 
}

//Generates a random key of URL_LENGTH and contains URL_CHARSET
func generateShortKey() string {
  shortKey := make([]byte, URL_LENGTH)

  for i := range shortKey {
    shortKey[i] = URL_CHARSET[rand.Intn(len(URL_CHARSET))]
  }

  return string(shortKey)
}
