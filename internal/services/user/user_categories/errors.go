package services

import "errors"

var (
	ErrCreateCategory      = errors.New("erro ao criar categoria")
	ErrFetchCategories     = errors.New("erro ao buscar categorias")
	ErrFetchCategory       = errors.New("erro ao buscar categoria")
	ErrUpdateCategory      = errors.New("erro ao atualizar categoria")
	ErrDeleteCategory      = errors.New("erro ao deletar categoria")
	ErrCategoryNotFound    = errors.New("categoria não encontrada")
	ErrInvalidCategoryName = errors.New("o nome da categoria é obrigatório")
	ErrInvalidCategory     = errors.New("categoria: objeto inválido")
	ErrCategoryIDRequired  = errors.New("categoria: ID da categoria é obrigatório")
	ErrCheckBeforeUpdate   = errors.New("erro ao verificar dados antes da atualização")
)
