package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var SERVER_PORT = os.Getenv("SERVER_PORT")
var SERVER_REDIRECT_URL = fmt.Sprintf("localhost:%s/sk/", SERVER_PORT)
var LOGGING_LEVEL = os.Getenv("LOGGING_LEVEL")
var dbConnection *DBConnection
var logger = setupLogger()

// Create a logger that outputs json logs
func setupLogger() *slog.Logger {
	if len(LOGGING_LEVEL) == 0 {
		LOGGING_LEVEL = "INFO"
	}

	// Convert the LOGGING_LEVEL to uppercase to make it case-insensitive
	level := strings.ToUpper(LOGGING_LEVEL)

	// Set the appropriate log level based on the environment variable
	var logLevel slog.Level
	switch level {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		// Default to INFO if the LOGGING_LEVEL is invalid
		logLevel = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel, // Set log level to DEBUG
	}))
}

// Open the connection to the database and save it as a global pointer variable
func handeDatabaseConnection() {
	// No error, so if we return, it has successfully connected
	dbConnect, err := connectToDatabase(logger)
	if err != nil {
		// Exit the program
		os.Exit(1)
	}
	dbConnection = dbConnect
}

// Serve the main index page
func indexPage(writer http.ResponseWriter, request *http.Request) {
	logger.Debug("Received HTTP request",
		slog.String("url", "/"),
		slog.String("Method", request.Method),
		slog.String("Address", request.RemoteAddr),
	)
	http.ServeFile(writer, request, "template/index.html")
}

// Handle the form submit in the page
func formSubmit(writer http.ResponseWriter, request *http.Request) {
	logger.Debug("Received HTTP request",
		slog.String("url", "/CreateShortUrl"),
		slog.String("Method", request.Method),
		slog.String("Address", request.RemoteAddr),
	)

	if request.Method == http.MethodPost {
		// Check that the request can be parsed
		err := request.ParseForm()
		if err != nil {
			logger.Debug("Unable to process form submit",
				slog.String("url", "/CreateShortUrl"),
				slog.String("Address", request.RemoteAddr),
				slog.Int("StatusCode", http.StatusBadRequest),
				slog.String("Error", err.Error()),
			)
			http.Error(writer, "Unable to process form", http.StatusBadRequest)
			return
		}

		oldurl := request.FormValue("enteredURL")

		//Valite the input to see if its in a form of a url
		valid, exists, err := validateIfURL(oldurl)
		if err != nil {
			logger.Error("Error validating url input",
				slog.String("enteredURL", oldurl),
				slog.String("Error", err.Error()),
				slog.Int("StatusCode", http.StatusInternalServerError),
			)
		}

		if !valid {
			logger.Debug("Input is not in form of url",
				slog.String("enteredURL", oldurl),
				slog.Int("StatusCode", http.StatusBadRequest),
			)
			http.Error(writer, "Input is not in form a of a url", http.StatusBadRequest)
			return
		}

		var shortKey string
		if exists {
			//Query the database to get the shortkey if one already exsits
			shortKey, err = dbConnection.findShortkeyUsingURL(oldurl)
			if err != nil {
				logger.Error("Error finding shortended url",
					slog.String("oldurl", oldurl),
					slog.String("Error", err.Error()),
					slog.Int("StatusCode", http.StatusInternalServerError),
				)
				http.Error(writer, "Error finding the shortened url", http.StatusInternalServerError)
				return
			}
		} else {
			shortKey, err = createShortKey() //Generate new shortkey

			if err != nil {
				logger.Error("Unable to generate shortkey",
					slog.String("Error", err.Error()),
				)
				http.Error(writer, "Unable to generate shortkey", http.StatusInternalServerError)
				return
			}

			//Add record to database
			err = dbConnection.addRecord(oldurl, shortKey)

			if err != nil {
				logger.Error("Unable to write record to database",
					slog.String("oldUrl", oldurl),
					slog.String("shortKey", shortKey),
					slog.String("Error", err.Error()),
					slog.Int("StatusCode", http.StatusInternalServerError),
				)
				http.Error(writer, "Unable to write record to database", http.StatusInternalServerError)
				return
			}
		}

		// Create a struct with the form data to pass to the template
		type FormData struct {
			GetURL string
			NewURL string
		}
		data := FormData{
			GetURL: SERVER_REDIRECT_URL,
			NewURL: shortKey,
		}

		// Open the newurl html file and use it as a template
		tmpl, err := template.ParseFiles("template/newurl.html")
		if err != nil {
			logger.Error("Unable to parse template",
				slog.String("template", "tempalte/newurl.html"),
				slog.String("Error", err.Error()),
				slog.Int("StatusCode", http.StatusInternalServerError),
			)
			http.Error(writer, "Unable to load template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(writer, data)
		if err != nil {
			logger.Error("Error executing template",
				slog.String("Error", err.Error()),
				slog.Int("StatusCode", http.StatusInternalServerError),
			)
			http.Error(writer, "Unable to render template", http.StatusInternalServerError)
			return
		}

		// Log that page loading was a success
		logger.Debug("Page rendering success",
			slog.String("url", "/CreateShortUrl"),
			slog.String("Method", request.Method),
			slog.String("Address", request.RemoteAddr),
		)
	} else {
		logger.Debug("Method not supported",
			slog.String("Mehod", request.Method),
			slog.String("url", "/CreateShortUrl"),
		)
		http.Error(writer, "Only POST method is supported", http.StatusMethodNotAllowed)
	}
}

// User provides a short key and redirect them to the corresponding url
func shortKeyHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract the shortkey from the URL
	parts := strings.Split(request.URL.Path, "/")
	if len(parts) < 3 {
		logger.Debug("Short Key not provided",
			slog.Int("StatusCode", http.StatusBadRequest),
		)
		http.Error(writer, "Shortkey not provided", http.StatusBadRequest)
		return
	}

	shortKey := parts[2] // Get the shortkey from the URL
	// Process the shortkey by searching the database

	oldurl, err := dbConnection.findURLUsingShortkey(shortKey)
	if err != nil {
		logger.Debug("Could not find corresponding url",
			slog.String("shortKey", shortKey),
			slog.String("Error", err.Error()),
			slog.Int("StatusCode", http.StatusBadRequest),
		)
		http.Error(writer, "Could not find corresponding url", http.StatusBadRequest)
		return
	}

	logger.Debug("Found corresponding url",
		slog.String("shortkey", shortKey),
		slog.String("url", oldurl),
	)

	//Redict the user to the new url
	http.Redirect(writer, request, oldurl, http.StatusFound)
}

func main() {
	logger.Info("Starting Server",
		slog.String("RequestedPort", SERVER_PORT),
	)

	handeDatabaseConnection()

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/CreateShortUrl", formSubmit)
	http.HandleFunc("/sk/", shortKeyHandler)

	// Start the HTTP server and log any errors using structured logging
	err := http.ListenAndServe(fmt.Sprintf(":%s", SERVER_PORT), nil)
	if err != nil {
		logger.Error("Failed to start HTTP server",
			slog.String("serverPort", SERVER_PORT),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
}
