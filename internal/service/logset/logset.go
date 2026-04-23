package logset

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

//go:generate mockgen -source=logset.go -destination=mock/mock_logset.go -package=mock

type Storage interface {
	Get(ctx context.Context, logsetID uuid.UUID) ([]map[string]any, error)
	Insert(ctx context.Context, logsetID uuid.UUID, logs []map[string]any) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Get(ctx context.Context, logsetID uuid.UUID) ([]map[string]any, error) {
	logs, err := s.storage.Get(ctx, logsetID)
	if err != nil {
		return nil, fmt.Errorf("get logs: %w", err)
	}

	return logs, nil
}

func (s *Service) Ingest(ctx context.Context, logs []map[string]any) (uuid.UUID, error) {
	logsetID := uuid.New()

	if err := s.storage.Insert(ctx, logsetID, logs); err != nil {
		return uuid.UUID{}, fmt.Errorf("insert logs: %w", err)
	}

	return logsetID, nil
}
