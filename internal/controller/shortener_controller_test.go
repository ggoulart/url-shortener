package controller

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenerController_ShortURL(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          string
		setup                func(*MockShortenerService)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "when failed to parse request body",
			requestBody:          "{",
			setup:                func(*MockShortenerService) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"failed parse request body: unexpected EOF"}`,
		},
		{
			name:        "when shortener service failed",
			requestBody: `{"longUrl": "https://bytebytego.com/courses/system-design-interview/design-a-url-shortener"}`,
			setup: func(m *MockShortenerService) {
				longURL := "https://bytebytego.com/courses/system-design-interview/design-a-url-shortener"
				m.On("Shortener", mock.AnythingOfType("*gin.Context"), longURL).Return(url.URL{}, errors.New("shortener service failed"))
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"error":"shortener service failed"}`,
		},
		{
			name:        "when successfuly shortens url",
			requestBody: `{"longUrl": "https://bytebytego.com/courses/system-design-interview/design-a-url-shortener"}`,
			setup: func(m *MockShortenerService) {
				longURL := "https://bytebytego.com/courses/system-design-interview/design-a-url-shortener"
				shortenURL, _ := url.Parse("https://gg.com/shorten")
				m.On("Shortener", mock.AnythingOfType("*gin.Context"), longURL).Return(*shortenURL, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"shortUrl":"https://gg.com/shorten"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockShortenerService{}
			tt.setup(m)

			c := NewShortenerController(m)

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = &http.Request{
				Body: io.NopCloser(strings.NewReader(tt.requestBody)),
			}

			c.ShortenURL(ctx)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedResponseBody, recorder.Body.String())
		})
	}
}

type MockShortenerService struct {
	mock.Mock
}

func (s *MockShortenerService) Retrieve(ctx context.Context, encodedKey string) (url.URL, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MockShortenerService) Shortener(ctx context.Context, shortURL string) (url.URL, error) {
	args := s.Called(ctx, shortURL)
	return args.Get(0).(url.URL), args.Error(1)
}
