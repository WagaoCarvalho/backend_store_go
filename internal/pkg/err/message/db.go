package err

import (
	"errors"
)

var (
	ErrDBConnURLNotDefined = errors.New("variável de ambiente DB_CONN_URL não definida")
	ErrDBParseConfig       = errors.New("erro ao parsear configuração do pool de conexão")
	ErrDBNewPool           = errors.New("erro ao criar novo pool de conexão")
	ErrDuplicate           = errors.New("erro já cadastrado")
	ErrDBPing              = errors.New("❌ banco de dados não iniciado")
)
