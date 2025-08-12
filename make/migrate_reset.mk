.PHONY: migrate_reset

migrate_reset:
	@echo "Resetando migrations (drop da tabela schema_migrations)..."
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "DROP TABLE IF EXISTS schema_migrations CASCADE;"
	@echo "Tabela schema_migrations removida com sucesso."
