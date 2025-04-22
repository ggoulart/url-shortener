package controller

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthController_Health(t *testing.T) {
	tests := []struct {
		name                 string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "when health service is successful",
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"postgres":true}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get("http://localhost:8080/api/v1/health")
			assert.NoError(t, err)

			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, tt.expectedResponseBody, string(body))
		})
	}
}
