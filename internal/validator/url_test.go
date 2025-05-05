package validate_test

import (
	"errors"
	"testing"

	validator "github.com/SamB032/Go-URL-Shortener/internal/validator"
	mocks "github.com/SamB032/Go-URL-Shortener/mocks"

	"github.com/stretchr/testify/assert"
)

func TestValidateURL(t *testing.T) {

	t.Run("Test with new valid URL", func(t *testing.T) {
		urlInput := "https://github.com/SamB032/Go-URL-Shortener"

		// Setup db mock
		dbMock := mocks.NewMockDbinterface(t)
		dbMock.On("CheckIfURLExists", urlInput).Return(false, nil)

		// Call ValidateURL
		valid, exists, err := validator.ValidateURL(urlInput, dbMock)

		// Assert outputs
		assert.True(t, valid, "URL should be valid")
		assert.False(t, exists, "URL should not already exist in database")
		assert.NoError(t, err, "No error should be raised")
	})

	t.Run("Test with existing URL", func(t *testing.T) {
		urlInput := "https://github.com/SamB032/Go-URL-Shortener"

		// Setup db mock
		dbMock := mocks.NewMockDbinterface(t)
		dbMock.On("CheckIfURLExists", urlInput).Return(true, nil)

		// Call ValidateURL
		valid, exists, err := validator.ValidateURL(urlInput, dbMock)

		// Assert Outputs
		assert.True(t, valid, "URL should be valid")
		assert.True(t, exists, "URL already exists in database")
		assert.NoError(t, err, "No error should be raised")
	})

	t.Run("Test with database error", func(t *testing.T) {
		urlInput := "https://github.com/SamB032/Go-URL-Shortener"

		// Setup db mock
		dbMock := mocks.NewMockDbinterface(t)
		dbMock.On("CheckIfURLExists", urlInput).Return(false, errors.New("database error"))

		// Call ValidateURL
		valid, exists, err := validator.ValidateURL(urlInput, dbMock)

		// Assert Outputs
		assert.False(t, valid, "URL should not be valid")
		assert.False(t, exists, "Exist in database should return false")
		assert.Error(t, err, "Error should be returned")
	})

	t.Run("Test with invalid url", func(t *testing.T) {
		urlInput := "invalid url"

		// Setup db mock
		dbMock := mocks.NewMockDbinterface(t)
		dbMock.On("CheckIfURLExists", urlInput).Return(false, nil)

		// Call ValidateURL
		valid, exists, err := validator.ValidateURL(urlInput, dbMock)

		// Assert Outputs
		assert.False(t, valid, "URL should not be valid")
		assert.False(t, exists, "URL should not exist in database")
		assert.NoError(t, err, "No error should be returend")
	})
}
