package storage

import (
	"context"
	"orchestrator/internal/core/enteties"
)

func (s *Storage) Login(ctx context.Context, user *enteties.User) (*enteties.User, error) {
	row := s.base.QueryRow(ctx,
		`
	SELECT id FROM users
	WHERE login = $1 AND password = $2;
	`,
		user.Login,
		user.Password,
	)
	if err := row.Scan(&user.ID); err != nil {
		return nil, err
	}
	if user.ID.String() == "" {
		return nil, enteties.ErrorLoginOrPasswordIncorrect
	}
	return user, nil
}

func (s *Storage) Authenticate(ctx context.Context, user *enteties.User) (*enteties.User, error) {
	return nil, nil
}
