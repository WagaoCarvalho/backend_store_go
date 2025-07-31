package services

import "errors"

var (
	ErrInvalidRelationData = errors.New("erro na validação dos dados")
	ErrRelationExists      = errors.New("erro relação já existe")
	ErrCreateRelation      = errors.New("erro ao criar relação")
	ErrGetRelations        = errors.New("erro ao buscar relações")
	ErrRelationNotFound    = errors.New("erro relação não encontrada")
	ErrDeleteRelation      = errors.New("erro ao deletar relação")
	ErrDeleteAllRelations  = errors.New("erro ao deletar todas as relações")
	ErrCheckRelationExists = errors.New("erro ao checar se relação existe")
)
