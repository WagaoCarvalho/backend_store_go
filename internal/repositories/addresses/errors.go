package repositories

import "errors"

var (
	ErrCreateAddress   = errors.New("erro ao criar endereço")
	ErrFetchAddress    = errors.New("erro ao buscar endereço")
	ErrAddressNotFound = errors.New("endereço não encontrado")
	ErrUpdateAddress   = errors.New("erro ao atualizar endereço")
	ErrDeleteAddress   = errors.New("erro ao excluir endereço")
)
