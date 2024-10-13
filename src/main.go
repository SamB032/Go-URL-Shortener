package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
)

const SERVER_REDIRECT_URL = "localhost:8080/sk/"

var dbConnection *sql.DB

// Open the connection to the database and save it as a global pointer variable
func handeDatabaseConnection() {
	// No error, so if we return, it has successfully connected
	message, db := connectToDatabase()
	log.Println(message)

	//Save db as a global variable
	dbConnection = db
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

		oldURL := request.FormValue("enteredURL")

		//Valite the input to see if its in a form of a url
		if !validateIfURL(oldURL){
			http.Error(writer, "Input is not in form a of a url", http.StatusBadRequest)
		}
		
		// Generate a new shortKey
		shortKey := generateShortKey() 		
		
		//Add record to database
		err = addRecord(oldURL, shortKey)
		if err != nil {
			log.Fatal(err)
			http.Error(writer, "Unable to write record to database", http.StatusInternalServerError)
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

		// Render the template with the form data
    err = tmpl.Execute(writer, data)
    if err != nil {
			log.Printf("Error executing template: %v", err)
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

	oldurl, err := findURLUsingShortkey(shortKey)
	if err != nil {
		log.Println(err)
		http.Error(writer, "Could not find corresponding url", http.StatusBadRequest)
	}

	//Redict the user to the new url
	http.Redirect(writer, request, oldurl, http.StatusFound)
}

func main() {
	log.Println("Server Starting")

	handeDatabaseConnection()

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/CreateShortUrl", formSubmit)
	http.HandleFunc("/sk/", shortKeyHandler)
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
