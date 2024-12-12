-- +goose Up
-- +goose StatementBegin
CREATE TYPE status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE IF NOT EXISTS gophermart_orders (
     order_id varchar(100) PRIMARY KEY,
     owner_user_id varchar(100) REFERENCES gophermart_users(user_id),
     status status NOT NULL,
     accrual DECIMAL,
     created_at TIMESTAMP DEFAULT (current_timestamp AT TIME ZONE 'Europe/Moscow')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_orders;
DROP TYPE status;
-- +goose StatementEnd
