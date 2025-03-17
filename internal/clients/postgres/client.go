package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Client struct {
	DB *sql.DB
}

func NewClient(config Config) (*Client, error) {
	db, err := sql.Open("postgres", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %v", err)
	}

	return &Client{DB: db}, nil
}

func (c *Client) Ping() error {
	err := c.DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping postgres: %v", err)
	}

	return nil
}
