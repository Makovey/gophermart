package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Makovey/gophermart/internal/config"
	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
)

const (
	errUniqueViolatesCode = "23505"
)

type repo struct {
	log  logger.Logger
	conn *pgx.Conn
}

func NewPostgresRepo(log logger.Logger, cfg config.Config) service.GophermartRepository {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURI())
	if err != nil {
		log.Error("unable to connect to database", "error", err.Error())
		panic(err)
	}

	return &repo{log: log, conn: conn}
}

func (r *repo) RegisterNewUser(ctx context.Context, user model.RegisterUser) error {
	fn := "postgresql.RegisterNewUser"

	_, err := r.conn.Exec(
		ctx,
		`INSERT INTO gophermart_users (user_id, login, password_hash) VALUES ($1, $2, $3)`,
		user.UserID,
		user.Login,
		user.PasswordHash,
	)
	if err != nil {
		r.log.Error(fmt.Sprintf("%s: failed to execute new user", fn), "error", err)
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == errUniqueViolatesCode {
			return service.ErrLoginIsAlreadyExist
		}

		return service.ErrExecStmt
	}

	r.log.Info(fmt.Sprintf("%s: inserted new user with login", fn), "login", user.Login)
	return nil
}

func (r *repo) LoginUser(ctx context.Context, login string) (model.RegisterUser, error) {
	fn := "postgresql.LoginUser"

	row := r.conn.QueryRow(
		ctx,
		`SELECT user_id, login, password_hash FROM gophermart_users WHERE login = $1`,
		login,
	)

	var user model.RegisterUser
	err := row.Scan(&user.UserID, &user.Login, &user.PasswordHash)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			r.log.Info(fmt.Sprintf("%s: user with login %s not found", fn, login))
			return model.RegisterUser{}, service.ErrNotFound
		default:
			r.log.Error(fmt.Sprintf("%s: failed to query user", fn), "error", err)
			return model.RegisterUser{}, service.ErrExecStmt
		}
	}

	return user, nil
}

func (r *repo) Close() error {
	return r.conn.Close(context.Background())
}
