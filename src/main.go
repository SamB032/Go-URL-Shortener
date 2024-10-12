package main

import (
	"fmt"
	"net/http"
	"log"
)

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

		// Return the data back to the client
		fmt.Fprintf(writer, "<html><body>")
    fmt.Fprintf(writer, "<h2>Form Submission Received</h2>")
    fmt.Fprintf(writer, "<p><strong>oldURL:</strong> %s</p>", oldURL)
    fmt.Fprintf(writer, "</body></html>")
	} else {
		http.Error(writer, "Only POST method is supported", http.StatusMethodNotAllowed)
	}
}

func main() {
	log.Println("Server Starting")

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/CreateShortUrl", formSubmit)
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}
