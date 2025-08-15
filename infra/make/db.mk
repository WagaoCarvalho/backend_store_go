.PHONY: db db_test stop_db clean_db migrate_up migrate_down

db:
	@echo "Subindo o banco de dados com Docker..."
	@docker compose --env-file .env up -d

db_test:
	@echo "Subindo o banco de dados para TESTES com Docker..."
	@docker compose -f docker-compose.test.yaml up -d

stop_db:
	@echo "Parando o banco de dados..."
	@docker-compose down

clean_db:
	@echo "Limpando containers e volumes..."
	@docker-compose down --volumes

migrate_up: migrate_up_user migrate_up_supplier migrate_up_client migrate_up_contact migrate_up_product migrate_up_address
	@echo "Todas as migrações foram aplicadas com sucesso."

migrate_down: migrate_down_user migrate_down_supplier migrate_down_client migrate_down_contact migrate_down_product migrate_down_address
	@echo "Todas as migrações foram aplicadas com sucesso."






