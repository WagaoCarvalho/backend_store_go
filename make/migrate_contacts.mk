.PHONY: migrate_create_contacts_table

migrate_create_contacts_table:
	@migrate create -ext sql -dir db/migrations -seq create_contacts_table
