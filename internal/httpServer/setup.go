package urlServer

import (
	"fmt"
	"log/slog"
	"net/http"

	database "github.com/SamB032/Go-URL-Shortener/internal/database"
)

type CreateShotKeyFunc func(dbConnection database.DBInterface) (string, error)
type ValidateUrlFunc func(url string, db database.DBInterface) (bool, bool, error)

type Server struct {
	logger         *slog.Logger
	database       database.DBInterface
	validateURL    ValidateUrlFunc
	createShortKey CreateShotKeyFunc
	redirectURL    string
	templatesDir   string
	mux            *http.ServeMux
}

func NewServer(serverPort string, logger *slog.Logger, db database.DBInterface, validateURL ValidateUrlFunc, createShortKey CreateShotKeyFunc, templatesDir string) *Server {
	mux := http.NewServeMux()

	server := &Server{
		logger:         logger,
		database:       db,
		validateURL:    validateURL,
		createShortKey: createShortKey,
		redirectURL:    fmt.Sprintf("localhost:%s/sk/", serverPort),
		templatesDir:   templatesDir,
		mux:            mux,
	}

	mux.HandleFunc("/", server.indexPage)
	mux.HandleFunc("/CreateShortUrl", server.formSubmit)
	mux.HandleFunc("/sk/", server.shortKeyHandler)

	return server
}

func (s *Server) Start(serverPort string) error {
	return http.ListenAndServe(fmt.Sprintf(":%s", serverPort), s.mux)
}
