# Nome do Projeto



## Tecnologias

- Go 1.23
- GitHub Actions para CI
- Codecov para cobertura de testes

## Estrutura do Projeto

- `/internal` - Código interno do serviço e modelos
- `/cmd` - Pontos de entrada da aplicação
- `/pkg` - Pacotes reutilizáveis

## Testes e Cobertura

- Os testes são executados usando `go test`.
- Cobertura de testes é enviada para Codecov via GitHub Actions.
- Para executar localmente:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
