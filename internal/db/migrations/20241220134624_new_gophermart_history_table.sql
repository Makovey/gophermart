-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gophermart_history (
   id SERIAL PRIMARY KEY,
   owner_user_id varchar(100) REFERENCES gophermart_users(user_id),
   order_id varchar(100) REFERENCES gophermart_orders(order_id),
   accrual DECIMAL,
   created_at TIMESTAMP DEFAULT (current_timestamp AT TIME ZONE 'Europe/Moscow')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_history;
-- +goose StatementEnd
