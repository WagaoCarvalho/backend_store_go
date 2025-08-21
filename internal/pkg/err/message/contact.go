package err

import (
	"errors"
)

var (
	ErrContactNotFound      = errors.New("contato não encontrado")
	ErrContactCreate        = errors.New("erro ao criar contato")
	ErrContactGet           = errors.New("erro ao buscar contato")
	ErrContactGetByUser     = errors.New("erro ao buscar contatos por user_id")
	ErrContactGetByClient   = errors.New("erro ao buscar contatos por client_id")
	ErrContactGetBySupplier = errors.New("erro ao buscar contatos por supplier_id")
	ErrContactScan          = errors.New("erro ao escanear contato")
	ErrContactUpdate        = errors.New("erro ao atualizar contato")
	ErrContactDelete        = errors.New("erro ao deletar contato")
	ErrContactInvalidID     = errors.New("erro o id tem que ser maior que 0")
	ErrUserContactsGet      = errors.New("erro ao buscar contatos do usuário")
	ErrClientContactsGet    = errors.New("erro ao buscar contatos do cliente")
	ErrSupplierContactsGet  = errors.New("erro ao buscar contatos do fornecedor")
	ErrContactID            = errors.New("erro id deve ser maior que 0")
)
