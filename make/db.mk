.PHONY: db db_test stop_db clean_db

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

migrate_up:
	@echo "Aplicando migrações..."
	@migrate -database ${DB_CONN_URL} -path db/migrations up

migrate_down:
	@echo "Revertendo migrações..."
	@migrate -database ${DB_CONN_URL} -path db/migrations down
