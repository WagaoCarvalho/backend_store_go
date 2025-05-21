# Carrega variáveis do .env
include .env
export $(shell sed 's/=.*//' .env)

# Inclui os módulos
include make/db.mk
include make/migrate_users.mk
include make/migrate_clients.mk
include make/migrate_suppliers.mk
include make/migrate_addresses.mk
include make/migrate_contacts.mk
include make/migrate_products.mk
include make/migrate_services.mk
include make/migrate_sales.mk
include make/migrate_purchases.mk
include make/migrate_misc.mk

.PHONY: server db db_test stop_db clean_db migrate_up migrate_down

server:
	@echo "Iniciando servidor Go..."
	@go run cmd/http/*.go
