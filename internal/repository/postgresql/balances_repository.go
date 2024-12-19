package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/service"
)

type balancesRepository struct {
	log  logger.Logger
	pool *pgxpool.Pool
}

func NewBalancesRepository(log logger.Logger, pool *pgxpool.Pool) service.BalancesRepository {
	return &balancesRepository{
		log:  log,
		pool: pool,
	}
}

func (b *balancesRepository) UpdateUsersBalance(ctx context.Context, userID string, reward decimal.Decimal) error {
	fn := "postgresql.GetOrderByID"

	_, err := b.pool.Exec(
		ctx,
		`INSERT INTO gophermart_balances (owner_user_id, accrual, updated_at) VALUES ($1, $2, $3) ON CONFLICT (owner_user_id)
		DO UPDATE SET accrual = gophermart_balances.accrual + excluded.accrual, updated_at = excluded.updated_at`,
		userID,
		reward,
		time.Now(),
	)
	if err != nil {
		b.log.Error(fmt.Sprintf("%s: failed to update users balance", fn), "error", err)
		return service.ErrExecStmt
	}

	return nil
}
