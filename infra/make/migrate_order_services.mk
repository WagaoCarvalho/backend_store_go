.PHONY: migrate_create_order_services_table migrate_create_order_service_categories_table migrate_create_order_service_category_relations_table migrate_up_order_service migrate_down_order_service

migrate_create_order_services_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_order_services_table

migrate_create_order_service_categories_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_order_service_categories_table

migrate_create_order_service_category_relations_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_order_service_category_relations_table

migrate_up_order_service:
	@echo "Aplicando migrações: order_service..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

migrate_down_order_service:
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
