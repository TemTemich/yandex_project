package services

import (
	"context"
	"fmt"
	"log"
	"orchestrator/internal/core/enteties"
	"regexp"
	"sync"

	"github.com/google/uuid"
)

type storage interface {
	AddExpression(ctx context.Context, expression *enteties.ArithmeticExpression) error
	GetExpression(ctx context.Context, expressionID uuid.UUID) (*enteties.ArithmeticExpression, error)
	UpdateOperation(ctx context.Context, id uuid.UUID, result string) error
	GetOperations(ctx context.Context) ([]*enteties.Operation, error)
	GetExpressions(ctx context.Context) ([]*enteties.ArithmeticExpression, error)
	CreateUser(ctx context.Context, user *enteties.User) error
	Login(context.Context, *enteties.User) (*enteties.User, error)
}

type calculator interface {
	Calculate(operation *enteties.Operation)
}

type Service struct {
	ctx          context.Context
	cancel       context.CancelFunc
	re           *regexp.Regexp
	operationsCh chan *enteties.ArithmeticExpression

	storage storage
	calc    calculator
}

func (s *Service) parseExpression(expression string) []*enteties.Operation {
	slOperation := s.re.FindAllString(expression, -1)
	operations := make([]*enteties.Operation, len(slOperation))
	for i, v := range slOperation {
		operations[i] = &enteties.Operation{
			ID: uuid.New(),
			Op: v,
		}
	}
	return operations
}

func (s *Service) AddExpression(ctx context.Context, expression string) (uuid.UUID, error) {
	idExpression := uuid.New()
	operations := s.parseExpression(expression)
	userID := ctx.Value("user_id").(string)
	expr := &enteties.ArithmeticExpression{
		ID:         idExpression,
		UserID:     userID,
		Operations: operations,
		Expression: expression,
	}
	select {
	case s.operationsCh <- expr:
		return idExpression, nil
	default:
		return uuid.Nil, fmt.Errorf("channel is full")
	}
}

func (s *Service) GetExpressions(ctx context.Context) ([]*enteties.ArithmeticExpression, error) {

	return s.storage.GetExpressions(ctx)
}

func (s *Service) GetExpression(ctx context.Context, id uuid.UUID) (*enteties.ArithmeticExpression, error) {
	return s.storage.GetExpression(ctx, id)
}

func (s *Service) UpdateOperation(id uuid.UUID, operation string) error {
	ctx := context.Background()
	if err := s.storage.UpdateOperation(ctx, id, operation); err != nil {
		log.Println("err update operation", id, operation)
		return err

	}
	log.Println("update operation", id, operation)
	return nil
}

func (s *Service) GetOperations(ctx context.Context) ([]*enteties.Operation, error) {
	return s.storage.GetOperations(ctx)
}

func (s *Service) Run() {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	for {
		select {
		case expr, ok := <-s.operationsCh:
			if !ok {
				s.cancel()
				return
			}
			s.storage.AddExpression(s.ctx, expr)
			wg := &sync.WaitGroup{}
			wg.Add(len(expr.Operations))
			for _, op := range expr.Operations {
				op := op
				go func() {
					defer wg.Done()
					s.calc.Calculate(op)
				}()
			}
			wg.Wait()
		case <-s.ctx.Done():
			return
		}
	}
}

func NewService(
	storage storage,
	calc calculator,
) *Service {
	patternString := `((?:^|\+|\-)\d+(?:[\*|\/]\d+)*)`
	re, _ := regexp.Compile(patternString)
	return &Service{
		re:           re,
		operationsCh: make(chan *enteties.ArithmeticExpression),
		storage:      storage,
		calc:         calc,
	}
}
