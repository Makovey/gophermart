-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gophermart_users (
    user_id VARCHAR(100) PRIMARY KEY,
    login VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT (current_timestamp AT TIME ZONE 'Europe/Moscow')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_users;
-- +goose StatementEnd