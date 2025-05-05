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
		{"exists", "apple", true, nil, true, false},
		{"exists", "bananas", true, nil, true, false},
		{"not exists", "xyz", false, nil, false, false},
		{"not exists", "-jf", false, nil, false, false},
		{"query error", "err", false, errors.New("fail"), false, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close() //nolint:errcheck
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
		{"exists", "http://abc", true, nil, true, false},
		{"exists", "http://xyz", true, nil, true, false},
		{"not exists", "http://abc", false, nil, false, false},
		{"not exists", "http://xyz", false, nil, false, false},
		{"query error", "http://fail", false, errors.New("fail"), false, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close() //nolint:errcheck
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

func TestAddRecord(t *testing.T) {
	cases := []struct {
		name        string
		url         string
		shortKey    string
		mockError   error
		expectError bool
	}{
		{"success", "http://abc", "apples", nil, false},
		{"success", "http://xyz", "bananas", nil, false},
		{"exec error", "http://fail", "err", errors.New("fail"), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close() //nolint:errcheck
			conn := &database.Connection{Connection: db}

			exec := mock.ExpectExec(`INSERT INTO url \(created_at, old_url, shortkey\) VALUES \(\$1, \$2, \$3\)`).WithArgs(sqlmock.AnyArg(), tc.url, tc.shortKey)

			if tc.mockError != nil {
				exec.WillReturnError(tc.mockError)
			} else {
				exec.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err := conn.AddRecord(tc.url, tc.shortKey)
			if (err != nil) != tc.expectError {
				t.Errorf("got error %v, expect error %v", err != nil, tc.expectError)
			}
		})
	}
}

func TestFindURLUsingShortkey(t *testing.T) {
	cases := []struct {
		name        string
		shortKey    string
		mockURL     string
		mockError   error
		expectError bool
	}{
		{"found", "apples", "http://abc", nil, false},
		{"found", "bananas", "http://xyz", nil, false},
		{"not found", "abc", "", errors.New("not found"), true},
		{"not found", "xyz", "", errors.New("not found"), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close() //nolint:errcheck
			conn := &database.Connection{Connection: db}

			if tc.mockError != nil {
				mock.ExpectQuery(`SELECT old_url FROM url WHERE shortkey = \$1`).WithArgs(tc.shortKey).WillReturnError(tc.mockError)
			} else {
				mock.ExpectQuery(`SELECT old_url FROM url WHERE shortkey = \$1`).WithArgs(tc.shortKey).WillReturnRows(sqlmock.NewRows([]string{"old_url"}).AddRow(tc.mockURL))
			}

			res, err := conn.FindURLUsingShortkey(tc.shortKey)
			if (err != nil) != tc.expectError || res != tc.mockURL {
				t.Errorf("got (%v, %v), want (%v, %v)", res, err != nil, tc.mockURL, tc.expectError)
			}
		})
	}
}

func TestFindShortkeyUsingURL(t *testing.T) {
	cases := []struct {
		name        string
		url         string
		mockKey     string
		mockError   error
		expectError bool
	}{
		{"found", "http://abc", "apples", nil, false},
		{"found", "http://xyz", "bananas", nil, false},
		{"not found", "http://abc", "", errors.New("not found"), true},
		{"not found", "http://xyz", "", errors.New("not found"), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close() //nolint:errcheck
			conn := &database.Connection{Connection: db}

			if tc.mockError != nil {
				mock.ExpectQuery(`SELECT shortkey FROM url WHERE old_url = \$1`).WithArgs(tc.url).WillReturnError(tc.mockError)
			} else {
				mock.ExpectQuery(`SELECT shortkey FROM url WHERE old_url = \$1`).WithArgs(tc.url).WillReturnRows(sqlmock.NewRows([]string{"shortkey"}).AddRow(tc.mockKey))
			}

			res, err := conn.FindShortkeyUsingURL(tc.url)
			if (err != nil) != tc.expectError || res != tc.mockKey {
				t.Errorf("got (%v, %v), want (%v, %v)", res, err != nil, tc.mockKey, tc.expectError)
			}
		})
	}
}
