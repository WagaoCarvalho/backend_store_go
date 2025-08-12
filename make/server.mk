.PHONY: server

server:
	@echo "Iniciando servidor Go..."
	@go run cmd/http/*.go
