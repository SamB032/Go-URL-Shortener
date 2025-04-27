package urlServer

import (
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	shortkey "github.com/SamB032/Go-URL-Shortener/internal/shortKey"
	validate "github.com/SamB032/Go-URL-Shortener/internal/validator"
)

// Serve the main index page
func (s *Server) indexPage(writer http.ResponseWriter, request *http.Request) {
	s.logger.Debug("Received HTTP request",
		slog.String("url", "/"),
		slog.String("Method", request.Method),
		slog.String("Address", request.RemoteAddr),
	)
	http.ServeFile(writer, request, s.templatesDir + "index.html")
}

// Handle the form submit in the page
func (s *Server) formSubmit(writer http.ResponseWriter, request *http.Request) {
	s.logger.Debug("Received HTTP request",
		slog.String("url", "/CreateShortUrl"),
		slog.String("Method", request.Method),
		slog.String("Address", request.RemoteAddr),
	)

	if request.Method == http.MethodPost {
		// Check that the request can be parsed
		err := request.ParseForm()
		if err != nil {
			s.logger.Debug("Unable to process form submit",
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
		valid, exists, err := validate.ValidateURL(oldurl, s.database)
		if err != nil {
			s.logger.Error("Error validating url input",
				slog.String("enteredURL", oldurl),
				slog.String("Error", err.Error()),
				slog.Int("StatusCode", http.StatusInternalServerError),
			)
		}

		if !valid {
			s.logger.Debug("Input is not in form of url",
				slog.String("enteredURL", oldurl),
				slog.Int("StatusCode", http.StatusBadRequest),
			)
			http.Error(writer, "Input is not in form a of a url", http.StatusBadRequest)
			return
		}

		var shortKey string
		if exists {
			//Query the database to get the shortkey if one already exsits
			shortKey, err = s.database.FindShortkeyUsingURL(oldurl)
			if err != nil {
				s.logger.Error("Error finding shortended url",
					slog.String("oldurl", oldurl),
					slog.String("Error", err.Error()),
					slog.Int("StatusCode", http.StatusInternalServerError),
				)
				http.Error(writer, "Error finding the shortened url", http.StatusInternalServerError)
				return
			}
		} else {
			shortKey, err = shortkey.CreateShortKey(s.database) //Generate new shortkey

			if err != nil {
				s.logger.Error("Unable to generate shortkey",
					slog.String("Error", err.Error()),
				)
				http.Error(writer, "Unable to generate shortkey", http.StatusInternalServerError)
				return
			}

			//Add record to database
			err = s.database.AddRecord(oldurl, shortKey)

			if err != nil {
				s.logger.Error("Unable to write record to database",
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
			GetURL: s.redirectURL,
			NewURL: shortKey,
		}

		// Open the newurl html file and use it as a template
		tmpl, err := template.ParseFiles(s.templatesDir + "newurl.html")
		if err != nil {
			s.logger.Error("Unable to parse template",
				slog.String("template", s.templatesDir + "newurl.html"),
				slog.String("Error", err.Error()),
				slog.Int("StatusCode", http.StatusInternalServerError),
			)
			http.Error(writer, "Unable to load template", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(writer, data)
		if err != nil {
			s.logger.Error("Error executing template",
				slog.String("Error", err.Error()),
				slog.Int("StatusCode", http.StatusInternalServerError),
			)
			http.Error(writer, "Unable to render template", http.StatusInternalServerError)
			return
		}

		// Log that page loading was a success
		s.logger.Debug("Page rendering success",
			slog.String("url", "/CreateShortUrl"),
			slog.String("Method", request.Method),
			slog.String("Address", request.RemoteAddr),
		)
	} else {
		s.logger.Debug("Method not supported",
			slog.String("Mehod", request.Method),
			slog.String("url", "/CreateShortUrl"),
		)
		http.Error(writer, "Only POST method is supported", http.StatusMethodNotAllowed)
	}
}

// User provides a short key and redirect them to the corresponding url
func (s *Server) shortKeyHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract the shortkey from the URL
	parts := strings.Split(request.URL.Path, "/")
	if len(parts) < 3 {
		s.logger.Debug("Short Key not provided",
			slog.Int("StatusCode", http.StatusBadRequest),
		)
		http.Error(writer, "Shortkey not provided", http.StatusBadRequest)
		return
	}

	shortKey := parts[2] // Get the shortkey from the URL
	// Process the shortkey by searching the database

	oldurl, err := s.database.FindURLUsingShortkey(shortKey)
	if err != nil {
		s.logger.Debug("Could not find corresponding url",
			slog.String("shortKey", shortKey),
			slog.String("Error", err.Error()),
			slog.Int("StatusCode", http.StatusBadRequest),
		)
		http.Error(writer, "Could not find corresponding url", http.StatusBadRequest)
		return
	}

	s.logger.Debug("Found corresponding url",
		slog.String("shortkey", shortKey),
		slog.String("url", oldurl),
	)

	//Redict the user to the new url
	http.Redirect(writer, request, oldurl, http.StatusFound)
}
