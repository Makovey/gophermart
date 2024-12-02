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