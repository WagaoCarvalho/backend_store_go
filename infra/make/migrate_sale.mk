.PHONY: migrate_create_sales_table migrate_up_sales migrate_down_sales

migrate_create_sales_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_sales_table

migrate_up_sales:
	@echo "Aplicando migrações: sales..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

migrate_down_sales:
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
