package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/shopspring/decimal"
)

func (r *Repo) RecordUsersWithdraw(ctx context.Context, userID, orderID string, amount decimal.Decimal) error {
	fn := "postgresql.RecordUsersWithdraw"

	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO gophermart_history (owner_user_id, order_id, withdraw, created_at) VALUES ($1, $2, $3, $4)`,
		userID,
		orderID,
		amount,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("[%s] failed to post history stamp: %w", fn, service.ErrExecStmt)
	}

	return nil
}

func (r *Repo) GetUsersHistory(ctx context.Context, userID string) ([]model.Withdraw, error) {
	fn := "postgresql.GetUsersHistory"

	rows, err := r.pool.Query(
		ctx,
		`SELECT order_id, withdraw, created_at FROM gophermart_history 
	  	WHERE owner_user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to query history withdraw: %w", fn, err)
	}
	defer rows.Close()

	var withdraws []model.Withdraw
	for rows.Next() {
		var withdraw model.Withdraw
		err = rows.Scan(&withdraw.OrderID, &withdraw.Withdraw, &withdraw.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("[%s] failed to query history withdraw: %w", fn, err)
		}
		withdraws = append(withdraws, withdraw)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("[%s] failed to iterate history withdraw: %w", fn, err)
	}

	return withdraws, nil
}
