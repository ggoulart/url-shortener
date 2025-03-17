package repository

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestShortenerRepository_FindEncodedKey(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		want    string
		wantErr error
	}{
		{
			name: "when db failed",
			setup: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`SELECT encoded_key FROM urls WHERE long_url = $1`)).
					WithArgs("a-long-url").
					WillReturnError(errors.New("db error"))
			},
			wantErr: errors.New("failed to find encoded key: db error"),
		},
		{
			name: "when db has no long url",
			setup: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`SELECT encoded_key FROM urls WHERE long_url = $1`)).
					WithArgs("a-long-url").
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "when successfully find encoded key",
			setup: func(s sqlmock.Sqlmock) {
				row := sqlmock.NewRows([]string{"encoded_key"}).AddRow("a-encoded-key")
				s.ExpectQuery(regexp.QuoteMeta(`SELECT encoded_key FROM urls WHERE long_url = $1`)).
					WithArgs("a-long-url").
					WillReturnRows(row)
			},
			want: "a-encoded-key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setup(dbMock)

			r := NewShortenerRepository(db)

			got, err := r.FindEncodedKey(context.Background(), "a-long-url")

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestShortenerRepository_FindLongURL(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		want    url.URL
		wantErr error
	}{
		{
			name: "when db failed",
			setup: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`SELECT long_url FROM urls WHERE encoded_key = $1`)).
					WithArgs("a-encoded-key").
					WillReturnError(errors.New("db error"))
			},
			wantErr: errors.New("failed to find longURL: db error"),
		},
		{
			name: "when db has invalid URL",
			setup: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`SELECT long_url FROM urls WHERE encoded_key = $1`)).
					WithArgs("a-encoded-key").
					WillReturnRows(sqlmock.NewRows([]string{"long_url"}).AddRow("://missing-scheme.com"))
			},
			wantErr: errors.New(`failed to parse longURL: parse "://missing-scheme.com": missing protocol scheme`),
		},
		{
			name: "when successfully find longURL",
			setup: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(`SELECT long_url FROM urls WHERE encoded_key = $1`)).
					WithArgs("a-encoded-key").
					WillReturnRows(sqlmock.NewRows([]string{"long_url"}).AddRow("http://valid-url.com"))
			},
			want: url.URL{Scheme: "http", Host: "valid-url.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setup(dbMock)

			r := NewShortenerRepository(db)

			got, err := r.FindLongURL(context.Background(), "a-encoded-key")

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestShortenerRepository_SaveURL(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		wantErr error
	}{
		{
			name: "when failed to insert url",
			setup: func(s sqlmock.Sqlmock) {
				s.ExpectExec(regexp.QuoteMeta(`INSERT INTO urls (encoded_key, long_url) VALUES ($1, $2)`)).
					WithArgs("a-encoded-key", "a-long-url").
					WillReturnError(errors.New("db error"))
			},
			wantErr: errors.New("failed to insert url: db error"),
		},
		{
			name: "when successfully save url",
			setup: func(s sqlmock.Sqlmock) {
				s.ExpectExec(regexp.QuoteMeta(`INSERT INTO urls (encoded_key, long_url) VALUES ($1, $2)`)).
					WithArgs("a-encoded-key", "a-long-url").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, dbMock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			tt.setup(dbMock)

			r := NewShortenerRepository(db)

			got := r.SaveURL(context.Background(), "a-encoded-key", "a-long-url")

			assert.Equal(t, tt.wantErr, got)
		})
	}
}
