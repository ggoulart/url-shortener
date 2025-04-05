package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
)

type ShortenerRepository interface {
	FindEncodedKey(ctx context.Context, longURL string) (string, error)
	FindLongURL(ctx context.Context, encodedKey string) (url.URL, error)
	SaveURL(ctx context.Context, shortURL string, longURL string) error
}

type ShortenerService struct {
	repository    ShortenerRepository
	shortenerHost string
	uuidGenerator func() string
}

func NewShortenerService(repository ShortenerRepository, shortenerHost string, uuidGenerator func() string) *ShortenerService {
	return &ShortenerService{repository: repository, shortenerHost: shortenerHost, uuidGenerator: uuidGenerator}
}

func (s *ShortenerService) Shortener(ctx context.Context, longURL string) (url.URL, error) {
	encodedKey, err := s.repository.FindEncodedKey(ctx, longURL)
	if err != nil {
		return url.URL{}, err
	}

	if encodedKey != "" {
		return s.buildShortURL(encodedKey)
	}

	id := s.uuidGenerator()
	encodedKey = base64.RawURLEncoding.EncodeToString([]byte(id))
	if len(encodedKey) > 7 {
		encodedKey = encodedKey[:7]
	}

	err = s.repository.SaveURL(ctx, encodedKey, longURL)
	if err != nil {
		return url.URL{}, err
	}

	return s.buildShortURL(encodedKey)
}

func (s *ShortenerService) Retrieve(ctx context.Context, encodedKey string) (url.URL, error) {
	longURL, err := s.repository.FindLongURL(ctx, encodedKey)
	if err != nil {
		return url.URL{}, err
	}

	return longURL, nil
}

func (s *ShortenerService) buildShortURL(encodedKey string) (url.URL, error) {
	shortURL, err := url.Parse(s.shortenerHost + "/" + encodedKey)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to build short URL: %v", err))
		return url.URL{}, errors.New("failed to build short URL")
	}

	return *shortURL, nil
}
