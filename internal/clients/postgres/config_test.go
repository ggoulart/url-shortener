package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDSN(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedDSN string
	}{
		{
			name:        "generates DSN correctly",
			config:      &Config{Host: "localhost", Port: "5432", User: "user", Pass: "password", DBName: "testdb", SSLMode: "disable"},
			expectedDSN: "host=localhost port=5432 user=user password=password dbname=testdb sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualDSN := tt.config.DSN()
			assert.Equal(t, tt.expectedDSN, actualDSN)
		})
	}
}
