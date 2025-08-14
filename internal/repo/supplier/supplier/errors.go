package repositories

import "errors"

var (
	ErrSupplierNotFound = errors.New("fornecedor não encontrado")
	ErrSupplierCreate   = errors.New("erro ao criar fornecedor")
	ErrSupplierUpdate   = errors.New("erro ao atualizar fornecedor")
	ErrSupplierDelete   = errors.New("erro ao deletar fornecedor")
	ErrSupplierDisable  = errors.New("erro ao desativar fornecedor")
	ErrSupplierEnable   = errors.New("erro ao ativar fornecedor")
	ErrSuppliersGet     = errors.New("erro ao listar fornecedores")
	ErrSupplierGet      = errors.New("erro ao buscar fornecedor por ID")
	ErrVersionConflict  = errors.New("conflito de versão: o endereço foi modificado por outra operação")
)
