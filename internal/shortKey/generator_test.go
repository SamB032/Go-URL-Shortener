package shortkey_test

import (
	"errors"
	"regexp"
	"testing"

	shortkey "github.com/SamB032/Go-URL-Shortener/internal/shortKey"
	mocks "github.com/SamB032/Go-URL-Shortener/mocks"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestCreateShoryKey(t *testing.T) {
	t.Run("Test CreateShortKey happy", func(t *testing.T) {
		// Setup mock with does not exist and no error
		dbMock := mocks.NewMockDbinterface(t)
		dbMock.On("CheckShortkeyExists", mock.AnythingOfType("string")).Return(false, nil).Once()

		// Call CreateShortKey
		shortkey, err := shortkey.CreateShortKey(dbMock)

		// assertions
		assert.NotEmpty(t, shortkey, "Shortkey should be returned")
		assert.NoError(t, err)
	})

	t.Run("Test CreateShortKey unhappy", func(t *testing.T) {
		// Setup mock that returns an error
		dbMock := mocks.NewMockDbinterface(t)
		dbMock.On("CheckShortkeyExists", mock.AnythingOfType("string")).Return(false, errors.New("database error")).Once()

		// Call CreateShortKey
		shortkey, err := shortkey.CreateShortKey(dbMock)

		// assertions
		assert.Empty(t, shortkey, "Shortkey should be generated")
		assert.Error(t, err)
	})
}

func TestGenerateShortKey(t *testing.T) {
	generatedShortKey := shortkey.GenerateShortKey()

	// Check length
	assert.Len(t, generatedShortKey, shortkey.URL_LENGTH)

	// Check for characters used 
	matched, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, generatedShortKey)
	assert.NoError(t, err)
	assert.True(t, matched, "string contains non-alphanumeric characters")
}
