package services

import "errors"

var (
	ErrInvalidAddressData = errors.New("address: dados do endereço inválidos")
	ErrAddressIDRequired  = errors.New("address: ID do endereço é obrigatório")
	ErrUpdateAddress      = errors.New("address: erro ao atualizar")
	ErrInvalidID          = errors.New("address: id inválido")
	ErrAddressNotFound    = errors.New("address: endereço não encontrado")
)
