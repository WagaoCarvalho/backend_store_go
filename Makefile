ifneq (,$(wildcard .env))
  include .env
  export
endif


# Inclui os módulos
include make/db.mk
include make/server.mk
include make/migrate_reset.mk
include make/migrate_users.mk
include make/migrate_clients.mk
include make/migrate_suppliers.mk
include make/migrate_addresses.mk
include make/migrate_contacts.mk
include make/migrate_products.mk
#include make/migrate_services_antt.mk

.PHONY: print-env
print-env:
	@echo DB_HOST=$(DB_HOST)
	@echo DB_PORT=$(DB_PORT)
	@echo DB_USER=$(DB_USER)
	@echo DB_PASSWORD=$(DB_PASSWORD)
	@echo DB_NAME=$(DB_NAME)