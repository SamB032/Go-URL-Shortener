package main

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var SERVER_PORT = os.Getenv("SERVER_PORT")
var SERVER_REDIRECT_URL = fmt.Sprintf("localhost:%s/sk/", SERVER_PORT)
var dbConnection *DBConnection
var logger = setupLogger()

// Create a logger that outputs json logs
func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// Open the connection to the database and save it as a global pointer variable
func handeDatabaseConnection() {
	// No error, so if we return, it has successfully connected
	dbConnection = connectToDatabase(logger)
}

// Serve the main index page
func indexPage(writer http.ResponseWriter, request *http.Request) {
	log.Printf("index: Received %s request from %s", request.Method, request.RemoteAddr)
	http.ServeFile(writer, request, "template/index.html")
}

// Handle the form submit in the page
func formSubmit(writer http.ResponseWriter, request *http.Request) {
	log.Printf("formSubmit: Received %s request from %s", request.Method, request.RemoteAddr)

	if request.Method == http.MethodPost {
		// Check that the request can be parsed
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, "Unable to process form", http.StatusBadRequest)
			return
		}

		oldurl := request.FormValue("enteredURL")

		//Valite the input to see if its in a form of a url
		valid, exists, err := validateIfURL(oldurl)
		if err != nil {
			log.Fatal(err)
			http.Error(writer, "There was an error when validating url input", http.StatusInternalServerError)
			return
		}

		if !valid{
			http.Error(writer, "Input is not in form a of a url", http.StatusBadRequest)
			return
		}
		
		var shortKey string
		if exists {
			//Query the database to get the shortkey if one already exsits
			shortKey, err = dbConnection.findShortkeyUsingURL(oldurl)
			if err != nil {
				log.Fatal(err)
				http.Error(writer, "Error finding the shortened url", http.StatusInternalServerError)
				return
			}
		} else {
			shortKey, err = createShortKey() //Generate new shortkey	

			if err != nil {
				log.Fatal(err)
				http.Error(writer, "Unable to generate shortkey", http.StatusInternalServerError)
				return
			}

			//Add record to database
			err = dbConnection.addRecord(oldurl, shortKey)

			if err != nil {
				log.Fatal(err)
				http.Error(writer, "Unable to write record to database", http.StatusInternalServerError)
				return
			}
		}

		// Create a struct with the form data to pass to the template
		type FormData struct { GetURL string; NewURL string }
    data := FormData{
			GetURL: SERVER_REDIRECT_URL,
			NewURL: shortKey,
    }

		// Open the newurl html file and use it as a template
		tmpl, err := template.ParseFiles("template/newurl.html")
    if err != nil {
			log.Printf("Error parsing template: %v", err)
      http.Error(writer, "Unable to load template", http.StatusInternalServerError)
      return
    }

    err = tmpl.Execute(writer, data)
    if err != nil {
			log.Fatalf("Error executing template: %v", err)
      http.Error(writer, "Unable to render template", http.StatusInternalServerError)
    }
	} else {
		http.Error(writer, "Only POST method is supported", http.StatusMethodNotAllowed)
	}
}

//User provides a short key and redirect them to the corresponding url
func shortKeyHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract the shortkey from the URL
  parts := strings.Split(request.URL.Path, "/")
  if len(parts) < 3 {
		http.Error(writer, "Shortkey not provided", http.StatusBadRequest)
    return
  }

  shortKey := parts[2] // Get the shortkey from the URL
  // Process the shortkey by searching the database

	oldurl, err := dbConnection.findURLUsingShortkey(shortKey)
	if err != nil {
		log.Println(err)
		http.Error(writer, "Could not find corresponding url", http.StatusBadRequest)
	}

	//Redict the user to the new url
	http.Redirect(writer, request, oldurl, http.StatusFound)
}

func main() {
	logger.Info("Starting Server")

	handeDatabaseConnection()

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/CreateShortUrl", formSubmit)
	http.HandleFunc("/sk/", shortKeyHandler)
	
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", SERVER_PORT), nil))
}
