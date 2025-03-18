# Carrega variáveis do .env
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: server db stop_db clean_db migrate_create_users_table migrate_create_product_sales_table migrate_create_product_purchase_items_table migrate_create_services_table migrate_up migrate_down

server:
	@echo "Iniciando servidor Go..."
	@go run cmd/http/*.go

db:
	@echo "Subindo o banco de dados com Docker..."
	@docker compose --env-file .env up -d

stop_db:
	@echo "Parando o banco de dados..."
	@docker-compose down

clean_db:
	@echo "Limpando containers e volumes..."
	@docker-compose down --volumes

###########
# Users
###########
# users
migrate_create_users_table:
	@migrate create -ext sql -dir db/migrations -seq create_users_table

# users_categories
migrate_create_users_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_user_categories_table

# user_category_relations
migrate_create_user_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_user_category_relations_table

###########
# Clients
###########
# clients
migrate_create_clients_table:
	@migrate create -ext sql -dir db/migrations -seq create_clients_table

# client_categories
migrate_create_client_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_client_categories_table

# clients_category_relations
migrate_create_client_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_clients_categories_relations_table

###########
# Suppliers
###########
# suppliers
migrate_create_suppliers_table:
	@migrate create -ext sql -dir db/migrations -seq create_suppliers_table

# supplier_categories
migrate_create_supplier_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_supplier_categories_table

# supplier_category_relations
migrate_create_supplier_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_supplier_category_relations_table

###########
# Addresses
###########
# addresses
migrate_create_addresses_table:
	@migrate create -ext sql -dir db/migrations -seq create_addresses_table
###########



###########
# Contacts
###########
# contacts
migrate_create_contacts_table:
	@migrate create -ext sql -dir db/migrations -seq create_contacts_table

###########
# Products
###########
# products
migrate_create_products_table:
	@migrate create -ext sql -dir db/migrations -seq create_products_table

# product_categories
migrate_create_product_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_product_categories_table

# product_category_relations
migrate_create_product_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_product_category_relations_table
###########

###########
# Services
###########
# services
migrate_create_services_table:
	@migrate create -ext sql -dir db/migrations -seq create_services_table

# service_categories
migrate_create_service_category_table:
	@migrate create -ext sql -dir db/migrations -seq create_service_category_table

# services_categories_rel
migrate_create_service_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_service_category_relations_table
##########

##########
# Sales
##########
# sales
migrate_create_sales_table:
	@migrate create -ext sql -dir db/migrations -seq create_sales_table

# sale_categories
migrate_create_sale_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_sale_categories_table

# sale_category_relations
migrate_create_sale_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_sale_category_relations_table

# sale_product_relations
migrate_create_sale_product_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_sale_product_relations_table

# sale_services
migrate_create_sale_services_table:
	@migrate create -ext sql -dir db/migrations -seq create_sale_services_table

##########
# Purchases
##########
# purchases
migrate_create_purchases_table:
	@migrate create -ext sql -dir db/migrations -seq create_purchases_table

# purchase_categories
migrate_create_purchase_categories_table:
	@migrate create -ext sql -dir db/migrations -seq create_purchase_categories_table

# purchase_category_relations
migrate_create_purchase_category_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_purchase_category_relations_table

# purchase_product_relations
migrate_create_purchase_product_relations_table:
	@migrate create -ext sql -dir db/migrations -seq create_purchase_product_relations_table

# purchase_services
migrate_create_purchase_services_table:
	@migrate create -ext sql -dir db/migrations -seq create_purchase_services_table

###


##########




# client_payments
migrate_create_client_payments_table:
	@migrate create -ext sql -dir db/migrations -seq create_client_payments_table


# product_suppliers_assignments
migrate_create_product_suppliers_assignments_table:
	@migrate create -ext sql -dir db/migrations -seq create_product_suppliers_assignments_table

# supplier_payments
migrate_create_supplier_payments_table:
	@migrate create -ext sql -dir db/migrations -seq create_supplier_payments_table



# price history
migrate_create_price_history_table:
	@migrate create -ext sql -dir db/migrations -seq create_price_history_table

# activity_log
migrate_create_activity_log_table:
	@migrate create -ext sql -dir db/migrations -seq create_activity_log_table





# migrations
migrate_up:
	@echo "Aplicando migrações..."
	@migrate -database ${DB_CONN_URL} -path db/migrations up

migrate_douwn:
	@echo "Revertendo migrações..."
	@migrate -database ${DB_CONN_URL} -path db/migrations down
