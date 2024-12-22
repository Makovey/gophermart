-- +goose Up
-- +goose StatementBegin
ALTER TABLE gophermart_history DROP CONSTRAINT gophermart_history_order_id_fkey;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE gophermart_history ADD CONSTRAINT gophermart_history_order_id_fkey;
-- +goose StatementEnd
