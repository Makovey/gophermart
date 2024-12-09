package service

import "errors"

var (
	ErrGeneratePass        = errors.New("failed to generate password")
	ErrPasswordDoesntMatch = errors.New("password doesn't match")
	ErrExecStmt            = errors.New("statement cannot execute")
	ErrLoginIsAlreadyExist = errors.New("login already exists")
	ErrNotFound            = errors.New("not found data in database")
)
