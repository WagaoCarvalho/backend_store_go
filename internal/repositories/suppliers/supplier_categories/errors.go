package repositories

import "errors"

var (
	ErrSupplierCategoryNotFound = errors.New("categoria de fornecedor não encontrada")
	ErrSupplierCategoryCreate   = errors.New("erro ao criar categoria")
	ErrSupplierCategoryGetAll   = errors.New("erro ao buscar categorias")
	ErrSupplierCategoryScanRow  = errors.New("erro ao ler dados da categoria")
	ErrSupplierCategoryUpdate   = errors.New("erro ao atualizar categoria")
	ErrSupplierCategoryDelete   = errors.New("erro ao deletar categoria")
)
