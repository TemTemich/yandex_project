package ports

import (
	"context"
	"orchestrator/internal/core/enteties"

	"github.com/google/uuid"
)

type Service interface {
	AddExpression(ctx context.Context, expression string) (uuid.UUID, error)
	UpdateOperation(uuid.UUID, string) error
	GetOperations(context.Context) ([]*enteties.Operation, error)
	GetExpression(ctx context.Context, id uuid.UUID) (*enteties.ArithmeticExpression, error)
	GetExpressions(context.Context) ([]*enteties.ArithmeticExpression, error)
	CreateUser(ctx context.Context, user *enteties.User) error
	Login(ctx context.Context, user *enteties.User) (*enteties.User, error)
}
