.PHONY: migrate_create_services_antt_table migrate_create_services_antt_categories_table migrate_create_services_antt_category_relations_table migrate_up_services_antt migrate_down_services_antt

migrate_create_services_antt_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_services_antts_table

migrate_create_services_antt_categories_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_services_antt_categories_table

migrate_create_services_antt_category_relations_table:
	@migrate create -ext sql -dir infra/db/migrations -seq create_services_antt_category_relations_table

migrate_up_services_antt:
	@echo "Aplicando migrações: services_antt..."
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations up

migrate_down_services_antt:
	@migrate -database ${DB_CONN_URL} -path infra/db/migrations down
