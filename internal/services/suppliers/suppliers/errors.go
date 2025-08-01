package services

import "errors"

var (
	ErrGetSupplier                  = errors.New("erro na validação dos dados")
	ErrGetSupplierVersion           = errors.New("erro ao obter a versão")
	ErrSupplierNameRequired         = errors.New("nome do fornecedor é obrigatório")
	ErrInvalidSupplierName          = errors.New("nome do fornecedor inválido")
	ErrSupplierCreateFailed         = errors.New("erro ao criar fornecedor")
	ErrSupplierNotFound             = errors.New("fornecedor não encontrado")
	ErrSupplierVersionConflict      = errors.New("conflito de versão ao atualizar o fornecedor")
	ErrSupplierUpdate               = errors.New("erro ao atualizar fornecedor")
	ErrSupplierVersionRequired      = errors.New("versão do fornecedor é obrigatória")
	ErrInvalidSupplierID            = errors.New("ID do fornecedor é inválido")
	ErrInvalidSupplierIDForDeletion = errors.New("ID inválido para deletar fornecedor")
	ErrDisableSupplier              = errors.New("erro ao desabilitar fornecedor")
	ErrEnableSupplier               = errors.New("erro ao habilitar fornecedor")
)
