package services

import "errors"

var (
	ErrInvalidEmail    = errors.New("email inválido")
	ErrCreateSupplier  = errors.New("erro ao criar fornecedor")
	ErrGetSupplier     = errors.New("erro ao buscar fornecedor")
	ErrGetVersion      = errors.New("erro ao buscar fornecedor")
	ErrUpdateSupplier  = errors.New("erro ao atualizar fornecedor")
	ErrDeleteSupplier  = errors.New("erro ao deletar fornecedor")
	ErrInvalidVersion  = errors.New("versão inválida")
	ErrGetAllSuppliers = errors.New("erro ao buscar fornecedors")
	ErrDisableSupplier = errors.New("erro ao desabilitar fornecedor")
	ErrEnableSupplier  = errors.New("erro ao habilitar fornecedor")
)
