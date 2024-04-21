package services

import (
	"context"
	"orchestrator/internal/core/enteties"
)

func (s *Service) Login(ctx context.Context, user *enteties.User) (*enteties.User, error) {
	return s.storage.Login(ctx, user)
}
