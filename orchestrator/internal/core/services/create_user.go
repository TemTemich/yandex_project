package services

import (
	"context"
	"orchestrator/internal/core/enteties"
)

func (s *Service) CreateUser(ctx context.Context, user *enteties.User) error {
	return s.storage.CreateUser(ctx, user)
}
