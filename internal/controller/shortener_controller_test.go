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
		expectedError        error
	}{
		{
			name:          "when failed to parse request body",
			requestBody:   "{",
			setup:         func(*MockShortenerService) {},
			expectedError: ErrBadRequest,
		},
		{
			name:        "when shortener service failed",
			requestBody: `{"longUrl": "https://bytebytego.com/courses/system-design-interview/design-a-url-shortener"}`,
			setup: func(m *MockShortenerService) {
				longURL := "https://bytebytego.com/courses/system-design-interview/design-a-url-shortener"
				m.On("Shortener", mock.AnythingOfType("*gin.Context"), longURL).Return(url.URL{}, errors.New("shortener service failed"))
			},
			expectedError: errors.New("shortener service failed"),
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

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError.Error(), ctx.Errors[len(ctx.Errors)-1].Error())
			} else {
				assert.Equal(t, tt.expectedStatusCode, recorder.Code)
				assert.Equal(t, tt.expectedResponseBody, recorder.Body.String())
			}
		})
	}
}

func TestShortenerController_RetrieveURL(t *testing.T) {
	tests := []struct {
		name                string
		setup               func(*MockShortenerService)
		expectedStatusCode  int
		expectedRedirectURL string
		expectedError       error
	}{
		{
			name: "when failed to retrieve url",
			setup: func(m *MockShortenerService) {
				m.On("Retrieve", mock.AnythingOfType("*gin.Context"), "NGVmMjk").Return(url.URL{}, errors.New("shortener service failed"))
			},
			expectedError: errors.New("shortener service failed"),
		},
		{
			name: "when successfully retrieves url",
			setup: func(m *MockShortenerService) {
				m.On("Retrieve", mock.AnythingOfType("*gin.Context"), "NGVmMjk").Return(url.URL{Host: "some-url"}, nil)
			},
			expectedStatusCode:  http.StatusFound,
			expectedRedirectURL: "//some-url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockShortenerService{}
			tt.setup(m)

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = &http.Request{}
			ctx.Params = gin.Params{{Key: "encodedKey", Value: "NGVmMjk"}}

			c := NewShortenerController(m)

			c.RetrieveURL(ctx)

			ctx.Writer.WriteHeaderNow()

			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError.Error(), ctx.Errors[len(ctx.Errors)-1].Error())
			} else {
				assert.Equal(t, tt.expectedStatusCode, recorder.Code)
				assert.Equal(t, tt.expectedRedirectURL, recorder.Header().Get("Location"))
			}
		})
	}
}

type MockShortenerService struct {
	mock.Mock
}

func (s *MockShortenerService) Retrieve(ctx context.Context, encodedKey string) (url.URL, error) {
	args := s.Called(ctx, encodedKey)
	return args.Get(0).(url.URL), args.Error(1)
}

func (s *MockShortenerService) Shortener(ctx context.Context, shortURL string) (url.URL, error) {
	args := s.Called(ctx, shortURL)
	return args.Get(0).(url.URL), args.Error(1)
}
