package err

import (
	"errors"
)

var (
	// crud
	ErrCreate                = errors.New("erro ao criar")
	ErrGet                   = errors.New("erro ao buscar")
	ErrNotFound              = errors.New("não encontrado")
	ErrIDNotFound            = errors.New("id não encontrado")
	ErrUpdate                = errors.New("erro ao atualizar")
	ErrDelete                = errors.New("erro ao deletar")
	ErrZeroID                = errors.New("erro ID deve ser maior que zero")
	ErrScan                  = errors.New("erro ao ler banco de dados")
	ErrIterate               = errors.New("erro ao iterar")
	ErrPercentInvalid        = errors.New("erro porcentagem inválida")
	ErrInvalidQuantity       = errors.New("quantidade inválida")
	ErrInvalidLimit          = errors.New("limite inválido")
	ErrInvalidOffset         = errors.New("offset inválido")
	ErrInvalidOrderField     = errors.New("ordem dos args inválida")
	ErrInvalidOrderDirection = errors.New("direção da ordem inválida")
	ErrInvalidDateRange      = errors.New("range de data inválido")

	// version
	ErrGetVersion      = errors.New("erro ao buscar versão")
	ErrVersionConflict = errors.New("erro de conflito de versão")
	ErrDisable         = errors.New("erro ao desabilitar")
	ErrEnable          = errors.New("erro ao ativar")

	ErrRelationExists = errors.New("relação já existe")
	ErrRelationCheck  = errors.New("erro ao verificar relação")

	ErrInvalidData = errors.New("erro na validação dos dados")

	ErrInvalidForeignKey = errors.New("chave estrangeira inválida")
)
