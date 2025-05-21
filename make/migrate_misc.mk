.PHONY: migrate_create_client_payments_table migrate_create_product_suppliers_assignments_table migrate_create_supplier_payments_table migrate_create_price_history_table migrate_create_activity_log_table

migrate_create_client_payments_table:
	@migrate create -ext sql -dir db/migrations -seq create_client_payments_table

migrate_create_product_suppliers_assignments_table:
	@migrate create -ext sql -dir db/migrations -seq create_product_suppliers_assignments_table

migrate_create_supplier_payments_table:
	@migrate create -ext sql -dir db/migrations -seq create_supplier_payments_table

migrate_create_price_history_table:
	@migrate create -ext sql -dir db/migrations -seq create_price_history_table

migrate_create_activity_log_table:
	@migrate create -ext sql -dir db/migrations -seq create_activity_log_table
