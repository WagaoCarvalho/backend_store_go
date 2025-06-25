package repositories

import "errors"

var (
	ErrRelationNotFound       = errors.New("relação usuário-categoria não encontrada")
	ErrRelationExists         = errors.New("relação já existe")
	ErrCreateRelation         = errors.New("erro ao criar relação")
	ErrGetRelationsByUser     = errors.New("erro ao buscar relações por usuário")
	ErrGetRelationsByCategory = errors.New("erro ao buscar relações por categoria")
	ErrScanRelation           = errors.New("erro ao ler relação")
	ErrIterateRelations       = errors.New("erro após ler relações")
	ErrDeleteRelation         = errors.New("erro ao deletar relação")
	ErrDeleteAllUserRelations = errors.New("erro ao deletar todas as relações do usuário")
	ErrInvalidForeignKey      = errors.New("usuário ou categoria inválido")
)
