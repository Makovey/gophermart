package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"

	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
)

func (r *Repo) IncreaseUsersBalance(ctx context.Context, userID string, reward decimal.Decimal) error {
	fn := "postgresql.IncreaseUsersBalance"

	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO gophermart_balances (owner_user_id, accrual, updated_at) VALUES ($1, $2, $3) ON CONFLICT (owner_user_id)
		DO UPDATE SET accrual = gophermart_balances.accrual + excluded.accrual, updated_at = $3`,
		userID,
		reward,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("[%s] failed to update users balance: %w", fn, service.ErrExecStmt)
	}

	return nil
}

func (r *Repo) DecreaseUsersBalance(ctx context.Context, userID string, withdraw decimal.Decimal) error {
	fn := "postgresql.DecreaseUsersBalance"

	tag, err := r.pool.Exec(
		ctx,
		`UPDATE gophermart_balances SET accrual = accrual - $1, withdrawn = withdrawn + $1, updated_at = $3 WHERE owner_user_id = $2`,
		withdraw,
		userID,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("[%s] failed to update users balance: %w", fn, service.ErrExecStmt)
	}

	if tag.RowsAffected() == 0 {
		return service.ErrNotFound
	}

	return nil
}

func (r *Repo) GetUsersBalance(ctx context.Context, userID string) (model.Balance, error) {
	fn := "postgresql.GetUsersBalance"

	row := r.pool.QueryRow(
		ctx,
		`SELECT accrual, withdrawn FROM gophermart_balances WHERE owner_user_id = $1`,
		userID,
	)
	var balance model.Balance
	err := row.Scan(&balance.Accrual, &balance.Withdrawn)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return model.Balance{}, service.ErrNotFound
		default:
			return model.Balance{}, fmt.Errorf("[%s] failed to query users balance: %w", fn, service.ErrExecStmt)
		}
	}

	return balance, nil
}
