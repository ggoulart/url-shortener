package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ggoulart/url-shortener/internal/controller"
	"github.com/ggoulart/url-shortener/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler_Unit(t *testing.T) {
	tests := []struct {
		name           string
		errToAttach    error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "bad request error",
			errToAttach:    controller.ErrBadRequest,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"` + controller.ErrBadRequest.Error() + `"}`,
		},
		{
			name:           "not found error",
			errToAttach:    repository.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"` + repository.ErrNotFound.Error() + `"}`,
		},
		{
			name:           "internal server error",
			errToAttach:    errors.New("something broke"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"something broke"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(resp)

			c.Error(tt.errToAttach)

			middlewareFunc := ErrorHandler()
			middlewareFunc(c)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.JSONEq(t, tt.expectedBody, resp.Body.String())
		})
	}
}
