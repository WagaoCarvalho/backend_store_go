package services

import "errors"

var (
	ErrCategoryNameRequired    = errors.New("nome da categoria é obrigatório")
	ErrCategoryIDInvalid       = errors.New("ID inválido")
	ErrCategoryIDRequired      = errors.New("ID da categoria é obrigatório")
	ErrCategoryDeleteInvalidID = errors.New("ID inválido para exclusão")
	ErrUpdateCategory          = errors.New("erro ao atualizar categoria")
	ErrCreateCategory          = errors.New("erro ao criar categoria")
	ErrGetCategory             = errors.New("erro ao buscar categoria")
	ErrGetCategories           = errors.New("erro ao buscar categorias")
	ErrDeleteCategory          = errors.New("erro ao deletar categorias")
)
