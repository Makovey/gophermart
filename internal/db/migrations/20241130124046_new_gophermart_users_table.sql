-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gophermart_users (
    user_id VARCHAR(100) PRIMARY KEY UNIQUE,
    login VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_users CASCADE;
-- +goose StatementEnd
