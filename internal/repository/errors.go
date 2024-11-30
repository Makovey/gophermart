package repository

import "errors"

var (
	ErrPrepareStmt         = errors.New("cannot prepare statement")
	ErrExecStmt            = errors.New("statement cannot execute")
	ErrLoginIsAlreadyExist = errors.New("login already exists")
)
