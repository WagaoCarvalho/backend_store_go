.PHONY: db db_test stop_db clean_db migrate_up migrate_down

# Subir banco de dados
db:
	@echo "Subindo o banco de dados com Docker..."
	@docker compose --env-file .env up -d

# Subir banco de dados para testes
db_test:
	@echo "Subindo o banco de dados para TESTES com Docker..."
	@docker compose -f docker-compose.test.yaml up -d

# Parar banco de dados
stop_db:
	@echo "Parando o banco de dados..."
	@docker compose down

# Limpar containers e volumes
clean_db:
	@echo "Limpando containers e volumes..."
	@docker compose down --volumes

migrate_up: migrate_up_user_all \
            migrate_up_supplier_all \
            migrate_up_client_all \
            migrate_up_contact \
            migrate_up_product \
            migrate_up_address
	@echo "Todas as migrações foram aplicadas com sucesso."

migrate_down: migrate_down_user_all \
              migrate_down_supplier_all \
              migrate_down_client_all \
              migrate_down_contact \
              migrate_down_product \
              migrate_down_address
	@echo "Todas as migrações foram revertidas com sucesso."
