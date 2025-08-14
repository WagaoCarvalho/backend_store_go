package repositories

import "errors"

var (
	ErrCategoryNotFound  = errors.New("categoria não encontrada")
	ErrGetCategories     = errors.New("erro ao buscar categorias")
	ErrScanCategory      = errors.New("erro ao ler os dados da categoria")
	ErrIterateCategories = errors.New("erro ao iterar sobre os resultados")
	ErrGetCategoryByID   = errors.New("erro ao buscar categoria por ID")
	ErrCreateCategory    = errors.New("erro ao criar categoria")
	ErrUpdateCategory    = errors.New("erro ao atualizar categoria")
	ErrDeleteCategory    = errors.New("erro ao deletar categoria")
)
