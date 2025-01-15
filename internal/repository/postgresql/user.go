package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
)

const (
	errUniqueViolatesCode = "23505"
)

func (r *Repo) RegisterNewUser(ctx context.Context, user model.RegisterUser) error {
	fn := "postgresql.RegisterNewUser"

	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO gophermart_users (user_id, login, password_hash) VALUES ($1, $2, $3)`,
		user.UserID,
		user.Login,
		user.PasswordHash,
	)
	if err != nil {
		r.log.Error(fmt.Sprintf("[%s] failed to execute new user", fn), "error", err)
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == errUniqueViolatesCode {
			return fmt.Errorf("[%s]: %w", fn, service.ErrLoginIsAlreadyExist)
		}

		return fmt.Errorf("[%s]: %w", fn, service.ErrExecStmt)
	}

	return nil
}

func (r *Repo) LoginUser(ctx context.Context, login string) (model.RegisterUser, error) {
	fn := "postgresql.LoginUser"

	row := r.pool.QueryRow(
		ctx,
		`SELECT user_id, login, password_hash FROM gophermart_users WHERE login = $1`,
		login,
	)

	var user model.RegisterUser
	err := row.Scan(&user.UserID, &user.Login, &user.PasswordHash)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return model.RegisterUser{}, fmt.Errorf("[%s]: %w", fn, service.ErrNotFound)
		default:
			return model.RegisterUser{}, fmt.Errorf("[%s]: %w", fn, service.ErrExecStmt)
		}
	}

	return user, nil
}
