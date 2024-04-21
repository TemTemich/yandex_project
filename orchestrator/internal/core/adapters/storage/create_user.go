package storage

import (
	"context"
	"orchestrator/internal/core/enteties"
)

func (s *Storage) CreateUser(ctx context.Context, user *enteties.User) error {
	row := s.base.QueryRow(
		ctx,
		`
		SELECT COUNT(*) FROM users
		WHERE login = $1;
		`,
		user.Login,
	)
	var count int
	if err := row.Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		return enteties.ErrorUserExist
	}

	if _, err := s.base.Exec(
		ctx, `
		INSERT INTO users (login, password, created_at) 
		VALUES ($1, $2, NOW())`,
		user.Login,
		user.Password,
	); err != nil {
		return err
	}

	return nil
}
