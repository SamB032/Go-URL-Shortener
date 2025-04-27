package database_test

import (
	"testing"
	"log/slog"

	database "github.com/SamB032/Go-URL-Shortener/internal/database"
	mocks "github.com/SamB032/Go-URL-Shortener/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DatabaseTestSuite struct {
	suite.Suite
	mockConn mocks.MockDbinterface
	logger   *slog.Logger
}

func TestConnectToDatabase_Happy(t *testing.T) {
	assert.True(t, false)
}
