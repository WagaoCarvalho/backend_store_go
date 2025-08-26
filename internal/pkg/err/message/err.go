package err

import (
	"errors"
)

var (
	// crud
	ErrCreate   = errors.New("erro ao criar")
	ErrGet      = errors.New("erro ao buscar")
	ErrNotFound = errors.New("não encontrado")
	ErrUpdate   = errors.New("erro ao atualizar")
	ErrDelete   = errors.New("erro ao deletar")
	ErrID       = errors.New("erro ID inválido")
	ErrScan     = errors.New("erro ao ler banco de dados")
	ErrIterate  = errors.New("erro ao iterar")

	// version
	ErrGetVersion      = errors.New("erro ao buscar versão")
	ErrVersionConflict = errors.New("erro de conflito de versão")
	ErrDisable         = errors.New("erro ao desabilitar")
	ErrEnable          = errors.New("erro ao ativar")

	ErrRelationExists = errors.New("relação já existe")
	ErrRelationCheck  = errors.New("erro ao verificar relação")

	ErrInvalidData = errors.New("erro na validação dos dados")
)
