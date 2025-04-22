package service

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenerService_Shortener(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*MockShortenerRepository)
		want    url.URL
		wantErr error
	}{
		{
			name: "when failed to findURL",
			setup: func(r *MockShortenerRepository) {
				r.On("FindEncodedKey", context.Background(), url.URL{Scheme: "http", Host: "some-long-url"}).Return("", errors.New("failed to find url"))
			},
			wantErr: errors.New("failed to find url"),
		},
		{
			name: "when found url failed to be build",
			setup: func(r *MockShortenerRepository) {
				r.On("FindEncodedKey", context.Background(), url.URL{Scheme: "http", Host: "some-long-url"}).Return("\x07", nil)
			},
			wantErr: errors.New("failed to build short URL"),
		},
		{
			name: "when successfully url already exists in db",
			setup: func(r *MockShortenerRepository) {
				r.On("FindEncodedKey", context.Background(), url.URL{Scheme: "http", Host: "some-long-url"}).Return("xZya7gG", nil)
			},
			want: url.URL{Scheme: "http", Host: "host-url.com", Path: "/api/v1/xZya7gG"},
		},
		{
			name: "when failed to save",
			setup: func(r *MockShortenerRepository) {
				r.On("FindEncodedKey", context.Background(), url.URL{Scheme: "http", Host: "some-long-url"}).Return("", nil)
				r.On("SaveURL", context.Background(), mock.Anything, url.URL{Scheme: "http", Host: "some-long-url"}).Return(errors.New("failed to save"))
			},
			wantErr: errors.New("failed to save"),
		},
		{
			name: "when successfully create shortURL and save it",
			setup: func(r *MockShortenerRepository) {
				r.On("FindEncodedKey", context.Background(), url.URL{Scheme: "http", Host: "some-long-url"}).Return("", nil)
				r.On("SaveURL", context.Background(), mock.Anything, url.URL{Scheme: "http", Host: "some-long-url"}).Return(nil)
			},
			want: url.URL{Scheme: "http", Host: "host-url.com", Path: "/api/v1/cmFuZG9"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MockShortenerRepository{}
			s := NewShortenerService(r, "http://host-url.com", func() string {
				return "random-generated-uuid"
			})
			tt.setup(r)

			got, err := s.Shortener(context.Background(), url.URL{Scheme: "http", Host: "some-long-url"})

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestShortenerService_Retrieve(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*MockShortenerRepository)
		want    url.URL
		wantErr error
	}{
		{
			name: "when failed to findURL",
			setup: func(r *MockShortenerRepository) {
				r.On("FindLongURL", context.Background(), "a-encoded-key").Return(url.URL{}, errors.New("failed to find url"))
			},
			wantErr: errors.New("failed to find url"),
		},
		{
			name: "when successfully findURL",
			setup: func(r *MockShortenerRepository) {
				r.On("FindLongURL", context.Background(), "a-encoded-key").Return(url.URL{Scheme: "http", Host: "host-url.com"}, nil)
			},
			want: url.URL{Scheme: "http", Host: "host-url.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MockShortenerRepository{}
			s := NewShortenerService(r, "http://host-url.com", func() string { return "" })
			tt.setup(r)

			got, err := s.Retrieve(context.Background(), "a-encoded-key")

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

type MockShortenerRepository struct {
	mock.Mock
}

func (m *MockShortenerRepository) FindEncodedKey(ctx context.Context, longURL url.URL) (string, error) {
	args := m.Called(ctx, longURL)
	return args.String(0), args.Error(1)
}

func (m *MockShortenerRepository) FindLongURL(ctx context.Context, encodedKey string) (url.URL, error) {
	args := m.Called(ctx, encodedKey)
	return args.Get(0).(url.URL), args.Error(1)
}

func (m *MockShortenerRepository) SaveURL(ctx context.Context, shortURL string, longURL url.URL) error {
	args := m.Called(ctx, shortURL, longURL)
	return args.Error(0)
}
