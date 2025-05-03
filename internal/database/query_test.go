package database_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SamB032/Go-URL-Shortener/internal/database"
)

func TestCheckShortkeyExists(t *testing.T) {
	cases := []struct {
		name         string
		shortKey     string
		mockExists   bool
		mockError    error
		expectResult bool
		expectError  bool
	}{
		{"exists", "abc", true, nil, true, false},
		{"not exists", "xyz", false, nil, false, false},
		{"query error", "err", false, errors.New("fail"), false, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			conn := &database.Connection{Connection: db}

			if tc.mockError != nil {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM url WHERE shortkey = \$1\)`).WithArgs(tc.shortKey).WillReturnError(tc.mockError)
			} else {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM url WHERE shortkey = \$1\)`).WithArgs(tc.shortKey).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(tc.mockExists))
			}

			res, err := conn.CheckShortkeyExists(tc.shortKey)
			if (err != nil) != tc.expectError || res != tc.expectResult {
				t.Errorf("got (%v, %v), want (%v, %v)", res, err != nil, tc.expectResult, tc.expectError)
			}
		})
	}
}

func TestCheckIfURLExists(t *testing.T) {
	cases := []struct {
		name         string
		url          string
		mockExists   bool
		mockError    error
		expectResult bool
		expectError  bool
	}{
		{"exists", "http://a", true, nil, true, false},
		{"not exists", "http://b", false, nil, false, false},
		{"query error", "http://fail", false, errors.New("fail"), false, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()
			conn := &database.Connection{Connection: db}

			if tc.mockError != nil {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 from URL where old_url = \$1\)`).WithArgs(tc.url).WillReturnError(tc.mockError)
			} else {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 from URL where old_url = \$1\)`).WithArgs(tc.url).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(tc.mockExists))
			}

			res, err := conn.CheckIfURLExists(tc.url)
			if (err != nil) != tc.expectError || res != tc.expectResult {
				t.Errorf("got (%v, %v), want (%v, %v)", res, err != nil, tc.expectResult, tc.expectError)
			}
		})
	}
}
