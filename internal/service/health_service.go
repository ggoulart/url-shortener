package service

import "context"

type DB interface {
	Ping() error
}

type HealthService struct {
	db DB
}

func NewHealthService(db DB) *HealthService {
	return &HealthService{db: db}
}

func (s *HealthService) Health(ctx context.Context) (map[string]bool, error) {
	m := map[string]bool{}

	err := s.db.Ping()
	m["postgres"] = err == nil

	return m, nil
}
