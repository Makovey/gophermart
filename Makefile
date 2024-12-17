include .env

SHELL := /bin/bash
LOCAL_MIGRATION_DIR=./internal/db/migrations
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/golang/mock/mockgen@v1.6.0

mig-s:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

mig-u:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

mig-d:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

test:
	go test ./...

remote-test:
	go build -o cmd/gophermart/gophermart cmd/gophermart/*.go
	./gophermarttest -test.v \
	-gophermart-binary-path=cmd/gophermart/gophermart \
	-gophermart-host=localhost \
	-gophermart-port=8080 \
	-gophermart-database-uri="postgres://admin:admin@localhost:5432/gophermart?sslmode=disable" \
	-accrual-binary-path=cmd/accrual/accrual_darwin_arm64 \
    -accrual-host=localhost \
    -accrual-port=8085 \
    -accrual-database-uri="postgres://admin:admin@localhost:5432/gophermart?sslmode=disable"