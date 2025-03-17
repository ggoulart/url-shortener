package postgres

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Host    string `mapstructure:"HOST"`
	Port    string `mapstructure:"PORT"`
	User    string `mapstructure:"USER"`
	Pass    string `mapstructure:"PASS"`
	DBName  string `mapstructure:"NAME"`
	SSLMode string `mapstructure:"SSL_MODE"`
}

func NewConfig() (*Config, error) {
	config := &Config{}
	err := viper.UnmarshalKey("db", config)
	if err != nil {
		return nil, fmt.Errorf("failed to load postgres config: %v", err)
	}

	return config, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.Host, c.Port, c.User, c.Pass, c.DBName, c.SSLMode)
}
