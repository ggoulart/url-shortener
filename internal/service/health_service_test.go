package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHealthService_Health(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*MockDBClient)
		want  map[string]bool
	}{
		{
			name: "when postgres is not healthy",
			setup: func(db *MockDBClient) {
				db.On("Ping").Return(errors.New("postgres is not healthy"))
			},
			want: map[string]bool{"postgres": false},
		},
		{
			name: "when postgres is healthy",
			setup: func(db *MockDBClient) {
				db.On("Ping").Return(nil)
			},
			want: map[string]bool{"postgres": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MockDBClient{}
			s := NewHealthService(m)
			tt.setup(m)

			got := s.Health(context.Background())

			assert.Equal(t, tt.want, got)
		})
	}
}

type MockDBClient struct {
	mock.Mock
}

func (m *MockDBClient) Ping() error {
	args := m.Called()
	return args.Error(0)
}
