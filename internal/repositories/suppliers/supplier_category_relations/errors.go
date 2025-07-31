package repositories

import "errors"

var (
	ErrRelationNotFound               = errors.New("relação supplier-categoria não encontrada")
	ErrCreateRelation                 = errors.New("erro ao criar relação")
	ErrCheckRelation                  = errors.New("erro ao verificar existência da relação")
	ErrGetRelationsBySupplier         = errors.New("erro ao buscar relações do fornecedor")
	ErrGetRelationsByCategory         = errors.New("erro ao buscar relações da categoria")
	ErrScanRelationRow                = errors.New("erro ao ler relação")
	ErrDeleteRelation                 = errors.New("erro ao deletar relação")
	ErrDeleteAllRelationsBySupplier   = errors.New("erro ao deletar todas as relações do fornecedor")
	ErrSupplierCategoryRelationUpdate = errors.New("erro ao atualizar a relação de categoria do fornecedor")
	ErrVersionConflict                = errors.New("conflito de versão: o registro foi modificado por outro processo")
	ErrSupplierNotFound               = errors.New("erro ao buscar fornecedor")
	ErrGetSuppliers                   = errors.New("erro ao buscar fornecedores")
	ErrSupplierUpdate                 = errors.New("erro ao atualizar fornecedor")
	ErrInvalidForeignKey              = errors.New("erro chave estrangeira inválida")
)
