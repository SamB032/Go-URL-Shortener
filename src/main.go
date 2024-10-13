package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

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
		type FormData struct { NewURL string }
    data := FormData{
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

func main() {
	log.Println("Server Starting")

	handeDatabaseConnection()

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/CreateShortUrl", formSubmit)
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
