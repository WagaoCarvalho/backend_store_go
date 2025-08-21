package err

import (
	"errors"
)

var (
	ErrInvalidForeignKey = errors.New("referência inválida: usuário, cliente ou fornecedor não existe")
	ErrDB                = errors.New("erro de banco de dados")
)
