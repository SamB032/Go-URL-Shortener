package url_server

import (
	"fmt"
	"log/slog"
	"net/http"

	database "github.com/SamB032/Go-URL-Shortener/internal/database"
)

type Server struct {
	logger       *slog.Logger
	database     database.DBInterface
	redirectURL  string
	templatesDir string
	mux          *http.ServeMux
}

func NewServer(serverPort string, logger *slog.Logger, db database.DBInterface, templatesDir string) *Server {
	mux := http.NewServeMux()

	server := &Server{
		logger:       logger,
		database:     db,
		redirectURL:  fmt.Sprintf("localhost:%s/sk/", serverPort),
		templatesDir: templatesDir,
		mux:          mux,
	}

	mux.HandleFunc("/", server.IndexPage)
	mux.HandleFunc("/CreateShortUrl", server.FormSubmit)
	mux.HandleFunc("/sk/", server.ShortKeyHandler)

	return server
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) Start(serverPort string) error {
	return http.ListenAndServe(fmt.Sprintf(":%s", serverPort), s.mux)
}
