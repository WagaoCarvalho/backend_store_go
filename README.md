
# Backend Store Go

Backend em Go para gerenciamento de produtos, fornecedores e clientes.

## Status do Projeto
[![Go Test Coverage](https://img.shields.io/badge/coverage-0%25-red)](https://codecov.io/gh/WagaoCarvalho/backend_store_go)

## Tecnologias
- Go 1.23
- PostgreSQL (opcional, dependendo do projeto)
- GitHub Actions para CI/CD
- Codecov para cobertura de testes
- Testify para testes unitários

## Estrutura do Projeto
```
cmd/        -> Aplicação principal
internal/   -> Pacotes internos (models, service, utils)
pkg/        -> Pacotes compartilháveis
tests/      -> Testes unitários adicionais
```
## Contribuição e Desenvolvimento Local

### Pré-requisitos
- Go >= 1.23
- Git
- Make (opcional)
- Dependências do projeto:
```bash
go mod tidy
```

### Rodando o projeto localmente
1. Clone o repositório:
```bash
git clone https://github.com/WagaoCarvalho/backend_store_go.git
cd backend_store_go
```

2. Configure variáveis de ambiente (exemplo `.env`):
```bash
cp .env.example .env
# Edite conforme necessário
```

3. Inicie o serviço (exemplo se tiver `main.go`):
```bash
go run cmd/main.go
```

### Testes
Para rodar todos os testes com saída no terminal e cobertura:

```bash
# Testes com relatório detalhado
go test -v -cover ./...

# Apenas cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Comandos Makefile (opcional)
Se tiver um `Makefile` configurado, você pode simplificar:

```bash
make deps     # Instala dependências
make test     # Roda todos os testes
make coverage # Roda testes e gera relatório de cobertura
```

### Padrões de qualidade
- Linting: `golangci-lint run`
- Vet: `go vet ./...`
- Cobertura de testes: `go test -cover ./...`
