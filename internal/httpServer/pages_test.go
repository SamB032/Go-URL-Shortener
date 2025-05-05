package url_server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	server "github.com/SamB032/Go-URL-Shortener/internal/httpServer"
	mocks "github.com/SamB032/Go-URL-Shortener/mocks"
)

func TestIndexPage(t *testing.T) {
	dbMock := mocks.NewMockDbinterface(t)

	server := server.NewServer("8080", logger, dbMock, TEMPLATES_DIR)
	assert.NotNil(t, server, "Expected server to be non-nil")

	// Create a request to the index page
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Call the handler
	server.IndexPage(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", resp.Header.Get("Content-Type"), "Unexpected Content type")
}

func TestFormSubmit(t *testing.T) {

	t.Run("Test FormSubmit_happy", func(t *testing.T) {
		// Assert with mock DB
		dbMock := mocks.NewMockDbinterface(t)
		dbMock.On("CheckIfURLExists", "https://example.com").Return(false, nil)
		dbMock.On("CheckShortkeyExists", mock.AnythingOfType("string")).Return(false, nil)
		dbMock.On("AddRecord", "https://example.com", mock.AnythingOfType("string")).Return(nil)

		// Create server
		server := server.NewServer("8080", logger, dbMock, TEMPLATES_DIR)
		assert.NotNil(t, server, "Expected server to be non-nil")

		// Setup a form submit resource
		form := url.Values{}
		form.Set("enteredURL", "https://example.com")
		req := httptest.NewRequest(http.MethodPost, "/CreateShortUrl", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp := httptest.NewRecorder()
		server.FormSubmit(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "text/html; charset=utf-8", resp.Header().Get("Content-Type"))
	})

	t.Run("Test FormSubmit_with_existing_url", func(t *testing.T) {
		// Assert with mock DB
		dbMock := mocks.NewMockDbinterface(t)
		dbMock.On("CheckIfURLExists", "https://example.com").Return(true, nil)
		dbMock.On("FindShortkeyUsingURL", mock.AnythingOfType("string")).Return("test123", nil)

		// Create server
		server := server.NewServer("8080", logger, dbMock, TEMPLATES_DIR)
		assert.NotNil(t, server, "Expected server to be non-nil")

		// Setup a form submit resource
		form := url.Values{}
		form.Set("enteredURL", "https://example.com")
		req := httptest.NewRequest(http.MethodPost, "/CreateShortUrl", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp := httptest.NewRecorder()
		server.FormSubmit(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "text/html; charset=utf-8", resp.Header().Get("Content-Type"))
		assert.Contains(t, resp.Body.String(), "test123")
	})

	t.Run("Test FormSubmit_with_db_error", func(t *testing.T) {
		// Assert with mock DB
		dbMock := mocks.NewMockDbinterface(t)
		dbMock.On("CheckIfURLExists", "https://example.com").Return(false, errors.New("DB error"))

		// Create server
		server := server.NewServer("8080", logger, dbMock, TEMPLATES_DIR)
		assert.NotNil(t, server, "Expected server to be non-nil")

		// Setup a form submit resource
		form := url.Values{}
		form.Set("enteredURL", "https://example.com")
		req := httptest.NewRequest(http.MethodPost, "/CreateShortUrl", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp := httptest.NewRecorder()
		server.FormSubmit(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Equal(t, "text/plain; charset=utf-8", resp.Header().Get("Content-Type"))
	})
}

func TestShortKeyHandler(t *testing.T) {

	t.Run("Test ShortKeyHandler_valid_key", func(t *testing.T) {
		// Mock DB with a valid URL for the short key
		dbMock := mocks.NewMockDbinterface(t)
		dbMock.On("FindURLUsingShortkey", "validkey").Return("https://example.com", nil)

		// Create server
		server := server.NewServer("8080", logger, dbMock, TEMPLATES_DIR)
		assert.NotNil(t, server, "Expected server to be non-nil")

		// Setup a valid short key request
		req := httptest.NewRequest(http.MethodGet, "/short/validkey", nil)
		resp := httptest.NewRecorder()
		server.ShortKeyHandler(resp, req)

		// Assert the response code and the redirect location
		assert.Equal(t, http.StatusFound, resp.Code)
		assert.Equal(t, "https://example.com", resp.Header().Get("Location"))
	})

	t.Run("Test ShortKeyHandler_missing_shortkey", func(t *testing.T) {
		// Create server with mock DB
		dbMock := mocks.NewMockDbinterface(t)
		server := server.NewServer("8080", logger, dbMock, TEMPLATES_DIR)

		// Simulate a missing short key in the URL
		req := httptest.NewRequest(http.MethodGet, "/short", nil)
		resp := httptest.NewRecorder()
		server.ShortKeyHandler(resp, req)

		// Assert the response status code and the content type
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, "text/plain; charset=utf-8", resp.Header().Get("Content-Type"))
	})
}
