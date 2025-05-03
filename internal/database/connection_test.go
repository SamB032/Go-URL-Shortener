package database_test

import (
	"testing"
	"os"
	"errors"
	"database/sql"
	"log/slog"

	database "github.com/SamB032/Go-URL-Shortener/internal/database"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewConnection_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	mock.ExpectPing().WillReturnError(nil)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	conn := database.NewConnection(mockDB, logger, "localhost", "testdb")
	if conn == nil {
		t.Fatal("expected non-nil connection")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestNewConnection_PingFails(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	mock.ExpectPing().WillReturnError(os.ErrPermission)

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	conn := database.NewConnection(mockDB, logger, "localhost", "testdb")
	if conn != nil {
		t.Fatal("expected nil connection when ping fails")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestConnectToDatabase_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()

	mock.ExpectPing()

	// Swap out sqlOpen for testing
	original := database.SqlOpen
	database.SqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return mockDB, nil
	}
	defer func() { database.SqlOpen = original }()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	conn := database.ConnectToDatabase("localhost", "5432", "user", "pass", "testdb", logger)
	if conn == nil {
		t.Fatal("expected non-nil connection")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet sqlmock expectations: %v", err)
	}
}

func TestConnectToDatabase_OpenFails(t *testing.T) {
	// Swap out sqlOpen to simulate failure
	original := database.SqlOpen
	database.SqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("open failed")
	}
	defer func() { database.SqlOpen = original }()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	conn := database.ConnectToDatabase("localhost", "5432", "user", "pass", "testdb", logger)
	if conn != nil {
		t.Fatal("expected nil connection when sql.Open fails")
	}
}
