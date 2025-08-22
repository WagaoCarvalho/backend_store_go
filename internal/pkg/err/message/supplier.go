package err

import (
	"errors"
)

var (
	ErrSupplierNotFound        = errors.New("fornecedor não encontrado")
	ErrSupplierCreate          = errors.New("erro ao criar fornecedor")
	ErrSupplierUpdate          = errors.New("erro ao atualizar fornecedor")
	ErrSupplierDelete          = errors.New("erro ao deletar fornecedor")
	ErrSupplierDisable         = errors.New("erro ao desativar fornecedor")
	ErrSupplierEnable          = errors.New("erro ao ativar fornecedor")
	ErrSuppliersGet            = errors.New("erro ao listar fornecedores")
	ErrSupplierGet             = errors.New("erro ao buscar fornecedor por ID")
	ErrSupplierIDInvalid       = errors.New("ID de fornecedor inválido")
	ErrSupplierVersionConflict = errors.New("conflito de versão: o endereço foi modificado por outra operação")

	ErrSupplierCategoryNotFound = errors.New("categoria de fornecedor não encontrada")
	ErrSupplierCategoryCreate   = errors.New("erro ao criar categoria")
	ErrSupplierCategoryGetAll   = errors.New("erro ao buscar categorias")
	ErrGetCategoryByID          = errors.New("erro ao buscar categoria pelo id")
	ErrSupplierCategoryScanRow  = errors.New("erro ao ler dados da categoria")
	ErrSupplierCategoryUpdate   = errors.New("erro ao atualizar categoria")
	ErrSupplierCategoryDelete   = errors.New("erro ao deletar categoria")
	ErrSupplierCategoryIterate  = errors.New("erro ao iterar sobre os resultados")
)
