package storage

import (
	"context"
	"log"
	"orchestrator/internal/core/enteties"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	base *pgxpool.Pool
}

func (s *Storage) GetExpressions(ctx context.Context) ([]*enteties.ArithmeticExpression, error) {
	rows, err := s.base.Query(
		ctx,
		`
		SELECT id, expression, status, result
		FROM expressions;
		`,
	)
	if err != nil {
		return nil, err
	}
	expressions := make([]*enteties.ArithmeticExpression, 0)
	for rows.Next() {
		var expr enteties.ArithmeticExpression
		if err := rows.Scan(&expr.ID, &expr.Expression, &expr.Status, &expr.Result); err != nil {
			return nil, err
		}
		expressions = append(expressions, &expr)
	}
	return expressions, nil
}

func (s *Storage) GetExpression(id uuid.UUID) (*enteties.ArithmeticExpression, error) {
	ctx := context.TODO()
	row, err := s.base.Query(ctx,
		`
	select status, result from expressions where id = $1;
	`, id)
	if err != nil {
		return nil, err
	}
	var expr enteties.ArithmeticExpression
	for row.Next() {
		if err := row.Scan(&expr.Status, &expr.Result); err != nil {
			log.Println(err)
			return nil, err
		}
		expr.ID = id
	}

	if expr.Status == "done" {
		return &expr, nil

	}

	row, err = s.base.Query(ctx, `
		SELECT SUM(sub.result::integer) as result
		FROM (
			SELECT expression_id, SUM(result::integer) as result
			FROM operations o
			GROUP BY expression_id
			HAVING COUNT(*) = COUNT(CASE WHEN status = 'done' THEN 1 END)
		) sub
		where sub.expression_id=$1;
	`, id)
	if err != nil {
		return nil, err

	}

	for row.Next() {
		if err := row.Scan(&expr.Result); err != nil {
			return nil, err
		}
	}
	if expr.Result == "" {
		return &expr, nil
	}

	if _, err := s.base.Exec(ctx,
		`UPDATE expressions SET status = 'done', done_at = current_timestamp, result=$2 WHERE id = $1;`,
		id, expr.Result); err != nil {
		return nil, err
	}
	expr.Status = "done"

	return &expr, nil
}

func (s *Storage) GetOperations(ctx context.Context) ([]*enteties.Operation, error) {
	rows, err := s.base.Query(ctx, `
		SELECT
			id,
			operation,
			result,
			EXTRACT(EPOCH FROM (done_at - created_at)) * 1000 AS duration_milliseconds
		FROM
			operations
		WHERE
			status = 'done';
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	operations := make([]*enteties.Operation, 0)
	for rows.Next() {
		operation := &enteties.Operation{}
		if err := rows.Scan(&operation.ID, &operation.Op, &operation.Result, &operation.Leadtime); err != nil {
			return nil, err

		}
		operations = append(operations, operation)
	}
	return operations, nil
}

func (s *Storage) AddExpression(expression *enteties.ArithmeticExpression) error {
	ctx := context.Background()
	tx, err := s.base.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Внести данные экспрешина
	if _, err := tx.Exec(ctx, `
		INSERT INTO expressions
		(id, "expression", created_at, done_at, status, "result")
		VALUES($1, $2, current_timestamp, null, 'work', '');
	`, expression.ID, expression.Expression); err != nil {
		return err
	}
	for _, operation := range expression.Operations {
		// Внести все операции
		if _, err := tx.Exec(ctx, `
			INSERT INTO operations
			(id, expression_id, operation, status, created_at, done_at, "result")
			VALUES($1, $2, $3, 'work', CURRENT_TIMESTAMP, null, null);
		`, operation.ID, expression.ID, operation.Op); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *Storage) UpdateOperation(ctx context.Context, id uuid.UUID, result string) error {
	tx, err := s.base.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		UPDATE operations
		SET status = 'done',
			done_at = CURRENT_TIMESTAMP,
			result = $2
		WHERE id = $1;
	`, id, result); err != nil {
		return err
	}
	return tx.Commit(ctx)

}

func NewStorage() *Storage {
	ctx := context.TODO()
	dsn := "postgres://serviceuser:service_password@storage:5432/service?sslmode=disable"
	base, _ := pgxpool.New(ctx, dsn)

	return &Storage{
		base: base,
	}
}
