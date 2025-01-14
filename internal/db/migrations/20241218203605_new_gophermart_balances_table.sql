-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gophermart_balances (
    id SERIAL PRIMARY KEY,
    owner_user_id varchar(100) UNIQUE REFERENCES gophermart_users(user_id),
    accrual DECIMAL,
    withdrawn DECIMAL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_balances;
-- +goose StatementEnd
