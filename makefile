-include .envrc
export

MIGRATIONS_PATH := ./cmd/migrate/migrations

.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(if $(filter-out $@,$(MAKECMDGOALS)),$(filter-out $@,$(MAKECMDGOALS)),1)

.PHONY: migrate-reset
migrate-reset:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down -all

.PHONY: seed
seed: 
	@go run ./cmd/migrate/seed/main.go  

.PHONY: seed-fresh
seed-fresh:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down -all
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up
	@go run ./cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g main.go -d cmd/api,internal/store && swag fmt
# Agar argument dinamis (seperti nama migrasi) tidak dianggap sebagai target Make
%:
	@: