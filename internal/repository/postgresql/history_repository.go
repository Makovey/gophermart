package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
)

type historyRepository struct {
	log  logger.Logger
	pool *pgxpool.Pool
}

func newHistoryRepository(log logger.Logger, pool *pgxpool.Pool) service.HistoryRepository {
	return &historyRepository{
		log:  log,
		pool: pool,
	}
}

func (h *historyRepository) RecordUsersWithdraw(ctx context.Context, userID, orderID string, amount decimal.Decimal) error {
	fn := "postgresql.RecordUsersWithdraw"

	_, err := h.pool.Exec(
		ctx,
		`INSERT INTO gophermart_history (owner_user_id, order_id, withdraw) VALUES ($1, $2, $3)`,
		userID,
		orderID,
		amount,
	)
	if err != nil {
		h.log.Error(fmt.Sprintf("%s: failed to post history stamp", fn), "error", err)
		return service.ErrExecStmt
	}

	return nil
}

func (h *historyRepository) GetUsersHistory(ctx context.Context, userID string) ([]model.Withdraw, error) {
	fn := "postgresql.GetUsersHistory"

	rows, err := h.pool.Query(
		ctx,
		`SELECT order_id, withdraw, created_at FROM gophermart_history 
	  	WHERE owner_user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		h.log.Error(fmt.Sprintf("%s: failed to query history withdraw", fn), "error", err)
		return nil, err
	}
	defer rows.Close()

	var withdraws []model.Withdraw
	for rows.Next() {
		var withdraw model.Withdraw
		err = rows.Scan(&withdraw.OrderID, &withdraw.Withdraw, &withdraw.CreatedAt)
		if err != nil {
			h.log.Error(fmt.Sprintf("%s: failed to scan history withdraw", fn), "error", err)
			return nil, err
		}
		withdraws = append(withdraws, withdraw)
	}

	if err = rows.Err(); err != nil {
		h.log.Error(fmt.Sprintf("%s: failed to iterate history withdraw", fn), "error", err)
		return nil, err
	}

	return withdraws, nil
}
