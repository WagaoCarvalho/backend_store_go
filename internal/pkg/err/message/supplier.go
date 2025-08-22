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
	ErrSupplierGetCategoryByID  = errors.New("erro ao buscar categoria pelo id")
	ErrSupplierCategoryScanRow  = errors.New("erro ao ler dados da categoria")
	ErrSupplierCategoryUpdate   = errors.New("erro ao atualizar categoria")
	ErrSupplierCategoryDelete   = errors.New("erro ao deletar categoria")
	ErrSupplierCategoryIterate  = errors.New("erro ao iterar sobre os resultados")

	ErrSupplierRelationExists         = errors.New("relação já existe")
	ErrSupplierRelationNotFound       = errors.New("relação supplier-categoria não encontrada")
	ErrSupplierCreateRelation         = errors.New("erro ao criar relação")
	ErrCheckRelation                  = errors.New("erro ao verificar existência da relação")
	ErrGetRelationsBySupplier         = errors.New("erro ao buscar relações do fornecedor")
	ErrSupplierGetRelationsByCategory = errors.New("erro ao buscar relações da categoria")
	ErrScanRelationRow                = errors.New("erro ao ler relação")
	ErrSupplierDeleteRelation         = errors.New("erro ao deletar relação")
	ErrDeleteAllRelationsBySupplier   = errors.New("erro ao deletar todas as relações do fornecedor")
	ErrSupplierCategoryRelationUpdate = errors.New("erro ao atualizar a relação de categoria do fornecedor")
	ErrVersionConflict                = errors.New("conflito de versão: o registro foi modificado por outro processo")
	ErrGetSuppliers                   = errors.New("erro ao buscar fornecedores")
)
