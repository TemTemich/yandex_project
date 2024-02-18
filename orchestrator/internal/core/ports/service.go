package ports

import (
	"context"
	"orchestrator/internal/core/enteties"

	"github.com/google/uuid"
)

type Service interface {
	AddExpression(expression string) (uuid.UUID, error)
	UpdateOperation(uuid.UUID, string) error
	GetOperations(context.Context) ([]*enteties.Operation, error)
	GetExpression(id uuid.UUID) (*enteties.ArithmeticExpression, error)
	GetExpressions(context.Context) ([]*enteties.ArithmeticExpression, error)
}
