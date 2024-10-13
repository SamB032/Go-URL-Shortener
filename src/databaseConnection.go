package main

import (
  _ "github.com/lib/pq"
  "github.com/joho/godotenv"

  "database/sql"
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

// TODO: Write to use docker-compose environment variables
//Load the posgres data from environment variables and return it as a struct
func getDatabaseInfo() PostgresData {
  err := godotenv.Load("../.env")
  // Check if the .env file was loaded correctly
  if err != nil {
    panic(err)
  }

  return PostgresData{
    os.Getenv("POSTGRES_HOST"), "5432", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"),
  }
}

func connectToDatabase() (string, *sql.DB) {
  dbInfo := getDatabaseInfo()

  psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
    dbInfo.host, dbInfo.port, dbInfo.user, dbInfo.password, dbInfo.dbName)

  //Initalise a connection to the database
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }

  // Make a ping to the database to see if its alive
  err = db.Ping()
  if err != nil {
    panic(err)
  }
  return "Successfully connected to the database", db
}
