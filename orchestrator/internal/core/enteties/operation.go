package enteties

import "github.com/google/uuid"

type Operation struct {
	ID       uuid.UUID
	Op       string
	Result   string
	Leadtime string
}

type ArithmeticExpression struct {
	ID         uuid.UUID
	UserID     string
	Expression string
	Status     string
	Operations []*Operation
	Result     string
}
