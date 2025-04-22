package controller

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenerController_ShortURL(t *testing.T) {
	tests := []struct {
		name                 string
		longUrl              string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "create short url",
			longUrl:              "https://dev.to/techschoolguru/load-config-from-file-environment-variables-in-golang-with-viper-2j2d",
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"shortUrl":"http://localhost:8080/api/v1/NGVmMjk"}`,
		},
		{
			name:                 "invalid body",
			longUrl:              "\")",
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error":"invalid body"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}

			reqBody := io.Reader(strings.NewReader(fmt.Sprintf(`{"longUrl": "%s"}`, tt.longUrl)))
			resp, err := client.Post("http://localhost:8080/api/v1/shorten", "application/json", reqBody)
			assert.NoError(t, err)

			respBody, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectedResponseBody, string(respBody))
		})
	}
}

func TestShortenerController_RetrieveURL(t *testing.T) {
	tests := []struct {
		name               string
		shortURL           string
		expectedStatusCode int
		expectedHeaders    http.Header
	}{
		{
			name:               "when url is found",
			shortURL:           "http://localhost:8080/api/v1/NGVmMjk",
			expectedStatusCode: http.StatusFound,
			expectedHeaders:    http.Header{"Location": []string{"https://dev.to/techschoolguru/load-config-from-file-environment-variables-in-golang-with-viper-2j2d"}},
		},
		{
			name:               "when url is not found",
			shortURL:           "http://localhost:8080/api/v1/312jnCa",
			expectedStatusCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			resp, err := client.Get(tt.shortURL)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectedHeaders.Get("Location"), resp.Header.Get("Location"))
		})
	}
}
