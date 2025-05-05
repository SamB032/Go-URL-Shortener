package database_test

import (
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	database "github.com/SamB032/Go-URL-Shortener/internal/database"
	"github.com/stretchr/testify/assert"
)

var logger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestNewConnection_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close() //nolint:errcheck

	mock.ExpectPing().WillReturnError(nil)

	conn, err := database.NewConnection(mockDB, logger, "localhost", "testdb")

	assert.NotNil(t, conn, "expected non-nil connection")
	assert.NoError(t, err, "No error should be raised")
}

func TestNewConnection_PingFails(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close() //nolint:errcheck

	mock.ExpectPing().WillReturnError(os.ErrPermission)

	conn, err := database.NewConnection(mockDB, logger, "localhost", "testdb")

	assert.Nil(t, conn, "Expected nil connection when sql.Open fails")
	assert.Error(t, err, "Expected error to be returned")
}

func TestConnectToDatabase_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close() //nolint:errcheck

	mock.ExpectPing()

	// Swap out sqlOpen for testing
	original := database.SqlOpen
	database.SqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return mockDB, nil
	}
	defer func() { database.SqlOpen = original }()

	conn, err := database.ConnectToDatabase("localhost", "5432", "user", "pass", "testdb", logger)

	assert.NotNil(t, conn, "expected non-nil connection")
	assert.NoError(t, err, "No error should be raised")
}

func TestConnectToDatabase_OpenFails(t *testing.T) {
	// Swap out sqlOpen to simulate failure
	original := database.SqlOpen
	database.SqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("open failed")
	}
	defer func() { database.SqlOpen = original }()

	conn, err := database.ConnectToDatabase("localhost", "5432", "user", "pass", "testdb", logger)

	assert.Nil(t, conn, "Expected nil connection when sql.Open fails")
	assert.Error(t, err, "Expected error to be returned")
}
