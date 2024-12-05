package repository

import "errors"

var (
	ErrExecStmt            = errors.New("statement cannot execute")
	ErrLoginIsAlreadyExist = errors.New("login already exists")
)
