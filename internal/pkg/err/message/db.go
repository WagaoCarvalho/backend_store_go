package err

import (
	"errors"
)

var (
	ErrInvalidForeignKey = errors.New("referência inválida: usuário, cliente ou fornecedor não existe")
	ErrDB                = errors.New("erro de banco de dados")

	ErrDBConnURLNotDefined = errors.New("variável de ambiente DB_CONN_URL não definida")
	ErrDBParseConfig       = errors.New("erro ao parsear configuração do pool de conexão")
	ErrDBNewPool           = errors.New("erro ao criar novo pool de conexão")
)
