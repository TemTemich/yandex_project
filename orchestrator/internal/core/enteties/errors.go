package enteties

import "errors"

var (
	ErrorUserExist                = errors.New("user already exist")
	ErrorLoginOrPasswordIncorrect = errors.New("login or password is incorrect")
)
