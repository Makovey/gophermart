-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gophermart_balances (
    id SERIAL PRIMARY KEY,
    owner_user_id varchar(100) REFERENCES gophermart_users(user_id),
    accrual DECIMAL,
    withdrawn DECIMAL,
    updated_at TIMESTAMP DEFAULT (current_timestamp AT TIME ZONE 'Europe/Moscow')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_balances;
-- +goose StatementEnd
