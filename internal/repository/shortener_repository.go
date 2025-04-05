package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
)

var ErrUnexpected = errors.New("unknown database error")
var ErrNotFound = errors.New("record not found")

type DB interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type ShortenerRepository struct {
	db DB
}

func NewShortenerRepository(db DB) *ShortenerRepository {
	return &ShortenerRepository{db: db}
}

func (r *ShortenerRepository) FindEncodedKey(ctx context.Context, longURL string) (string, error) {
	query := `SELECT encoded_key FROM urls WHERE long_url = $1`

	var encodedKey string
	err := r.db.QueryRowContext(ctx, query, longURL).Scan(&encodedKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		slog.Error(fmt.Sprintf("failed to find encoded key: %v", err))
		return "", ErrUnexpected
	}

	return encodedKey, nil
}

func (r *ShortenerRepository) FindLongURL(ctx context.Context, encodedKey string) (url.URL, error) {
	query := `SELECT long_url FROM urls WHERE encoded_key = $1`

	var dbLongURL string
	err := r.db.QueryRowContext(ctx, query, encodedKey).Scan(&dbLongURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn(fmt.Sprintf("encoded key %s not found: %v", encodedKey, err))
			return url.URL{}, ErrNotFound
		}

		slog.Error(fmt.Sprintf("failed to find longURL: %v", err))
		return url.URL{}, ErrUnexpected
	}

	longURL, err := url.Parse(dbLongURL)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to parse longURL: %v", err))
		return url.URL{}, ErrUnexpected
	}

	return *longURL, nil
}

func (r *ShortenerRepository) SaveURL(ctx context.Context, encodedKey string, longURL string) error {
	query := `INSERT INTO urls (encoded_key, long_url) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, encodedKey, longURL)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to insert url: %v", err))
		return ErrUnexpected
	}

	return nil
}
