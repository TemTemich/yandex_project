package enteties

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID
	Login    string `json:"login"`
	Password string `json:"password"`
}
