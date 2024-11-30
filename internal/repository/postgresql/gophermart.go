package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/Makovey/gophermart/internal/config"
	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/repository"
	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
)

const (
	insertNewUser = "insertNewUser"
)

type repo struct {
	log  logger.Logger
	conn *pgx.Conn
}

func (r *repo) RegisterNewUser(ctx context.Context, user model.RegisterUser) error {
	fn := "postgresql.RegisterNewUser"

	_, err := r.conn.Prepare(ctx, insertNewUser, `
		INSERT INTO gophermart_users (user_id, login, password_hash) VALUES ($1, $2, $3)
    `)
	if err != nil {
		r.log.Error(fmt.Sprintf("%s: failed to prepare new user", fn), "error", err)
		return repository.ErrPrepareStmt
	}

	_, err = r.conn.Exec(ctx, insertNewUser, user.UserID, user.Login, user.PasswordHash)
	if err != nil {
		r.log.Error(fmt.Sprintf("%s: failed to execute new user", fn), "error", err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return repository.ErrLoginIsAlreadyExist
		}

		return repository.ErrExecStmt
	}

	r.log.Info(fmt.Sprintf("%s: inserted new user with login", fn), "login", user.Login)
	return nil
}

func (r *repo) Close() error {
	return r.conn.Close(context.Background())
}

func NewPostgresRepo(log logger.Logger, cfg config.Config) service.GophermartRepository {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURI())
	if err != nil {
		log.Error("unable to connect to database", "error", err.Error())
		panic(err)
	}

	return &repo{log: log, conn: conn}
}
