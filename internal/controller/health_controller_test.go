package controller

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHealthController_Health(t *testing.T) {
	tests := []struct {
		name                 string
		setup                func(*MockHealthService)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "when health service is successful",
			setup: func(m *MockHealthService) {
				m.On("Health", mock.Anything, mock.Anything).Return(map[string]bool{"health": true})
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"health":true}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockHealthService{}
			tt.setup(m)
			h := NewHealthController(m)

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = &http.Request{}

			h.Health(ctx)

			assert.Equal(t, tt.expectedStatusCode, recorder.Code)
			assert.Equal(t, tt.expectedResponseBody, recorder.Body.String())
		})
	}
}

type MockHealthService struct {
	mock.Mock
}

func (h *MockHealthService) Health(ctx context.Context) map[string]bool {
	args := h.Called(ctx)
	return args.Get(0).(map[string]bool)
}
