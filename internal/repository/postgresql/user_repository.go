package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
)

const (
	errUniqueViolatesCode = "23505"
)

type userRepository struct {
	log  logger.Logger
	pool *pgxpool.Pool
}

func newUserRepository(log logger.Logger, pool *pgxpool.Pool) service.UserRepository {
	return &userRepository{
		log:  log,
		pool: pool,
	}
}

func (u *userRepository) RegisterNewUser(ctx context.Context, user model.RegisterUser) error {
	fn := "postgresql.RegisterNewUser"

	_, err := u.pool.Exec(
		ctx,
		`INSERT INTO gophermart_users (user_id, login, password_hash, created_at) VALUES ($1, $2, $3, $4)`,
		user.UserID,
		user.Login,
		user.PasswordHash,
		time.Now(),
	)
	if err != nil {
		u.log.Error(fmt.Sprintf("%s: failed to execute new user", fn), "error", err)
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == errUniqueViolatesCode {
			return service.ErrLoginIsAlreadyExist
		}

		return service.ErrExecStmt
	}

	return nil
}

func (u *userRepository) LoginUser(ctx context.Context, login string) (model.RegisterUser, error) {
	fn := "postgresql.LoginUser"

	row := u.pool.QueryRow(
		ctx,
		`SELECT user_id, login, password_hash FROM gophermart_users WHERE login = $1`,
		login,
	)

	var user model.RegisterUser
	err := row.Scan(&user.UserID, &user.Login, &user.PasswordHash)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			u.log.Info(fmt.Sprintf("%s: user with login %s not found", fn, login))
			return model.RegisterUser{}, service.ErrNotFound
		default:
			u.log.Error(fmt.Sprintf("%s: failed to query user", fn), "error", err)
			return model.RegisterUser{}, service.ErrExecStmt
		}
	}

	return user, nil
}
