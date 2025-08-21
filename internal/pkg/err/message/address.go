package err

import (
	"errors"
)

var (
	ErrAddressCreate   = errors.New("erro ao criar endereço")
	ErrAddressGet      = errors.New("erro ao buscar endereço")
	ErrAddressNotFound = errors.New("endereço não encontrado")
	ErrAddressUpdate   = errors.New("erro ao atualizar endereço")
	ErrAddressDelete   = errors.New("erro ao deletar endereço")
	ErrAddressID       = errors.New("erro id deve ser maior que 0")
)
