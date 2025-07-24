package services

import "errors"

var (
	ErrInvalidUserID              = errors.New("ID do usuário inválido")
	ErrInvalidCategoryID          = errors.New("ID da categoria inválido")
	ErrCreateRelation             = errors.New("erro ao criar relação")
	ErrCheckExistingRelation      = errors.New("erro ao verificar relação existente")
	ErrFetchUserRelations         = errors.New("erro ao buscar relações do usuário")
	ErrFetchCategoryRelations     = errors.New("erro ao buscar relações da categoria")
	ErrDeleteRelation             = errors.New("erro ao deletar relação")
	ErrDeleteAllUserRelations     = errors.New("erro ao deletar todas as relações do usuário")
	ErrCheckUserCategoryRelations = errors.New("erro ao verificar relações do usuário")
	ErrCheckRelationExists        = errors.New("erro ao verificar se a relação entre usuário e categoria existe")
)
