.PHONY: migrate_create_services_table migrate_create_service_category_table migrate_create_service_category_relations_table

migrate_create_services_table:
	@migrate create -ext sql -dir db/migrations -seq create_services_table

migrate_create_service_category_table:
	@migrate create -ext sql -dir db/migrations -seq create_service_category_table

migrate_create_service_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_service_category_relations_table
