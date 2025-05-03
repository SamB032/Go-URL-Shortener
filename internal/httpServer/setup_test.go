package url_server_test

import (
	"log/slog"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	assert "github.com/stretchr/testify/assert"

	server "github.com/SamB032/Go-URL-Shortener/internal/httpServer"
	mocks "github.com/SamB032/Go-URL-Shortener/mocks"
)

const TEMPLATES_DIR string = "../../templates/"

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestNewSever(t *testing.T) {
	dbMock := mocks.NewMockDbinterface(t)

	server := server.NewServer("8080", logger, dbMock, TEMPLATES_DIR)

	assert.NotNil(t, server, "Expected server to be non-nil")

	// Check if a handler is registered for "/"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, req)

	assert.Equal(t, rec.Code, http.StatusOK, "Expected 200 OK")
}

func TestStartServer(t *testing.T) {
	t.Run("Test StartServer Happy", func(t *testing.T) {
		dbMock := mocks.NewMockDbinterface(t)
		server := server.NewServer("8080", logger, dbMock, TEMPLATES_DIR)

		go func() {
			err := server.Start("8080")
			assert.NotNil(t, err, "There should be no error when starting the server")
		}()

		time.Sleep(1 * time.Second)
	})

	t.Run("Test StartServer Sad", func(t *testing.T) {
		dbMock := mocks.NewMockDbinterface(t)
		server := server.NewServer("80", logger, dbMock, TEMPLATES_DIR)

		err := server.Start("80")
		assert.Error(t, err, "StartServer should return an error on privileged port 80")
	})
}
