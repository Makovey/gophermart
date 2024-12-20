package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/Makovey/gophermart/internal/logger"
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
		`INSERT INTO gophermart_history (owner_user_id, order_id, accrual) VALUES ($1, $2, $3)`,
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
