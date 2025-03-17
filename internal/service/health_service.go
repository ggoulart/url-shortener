package service

import "context"

type DBClient interface {
	Ping() error
}

type HealthService struct {
	dbClient DBClient
}

func NewHealthService(dbClient DBClient) *HealthService {
	return &HealthService{dbClient: dbClient}
}

func (s *HealthService) Health(ctx context.Context) map[string]bool {
	m := map[string]bool{}

	err := s.dbClient.Ping()
	m["postgres"] = err == nil

	return m
}
