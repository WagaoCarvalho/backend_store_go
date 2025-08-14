package services

import "errors"

var (
	ErrContactNameRequired        = errors.New("nome do contato é obrigatório")
	ErrContactAssociationRequired = errors.New("o contato deve estar associado a um usuário, cliente ou fornecedor")
	ErrInvalidEmail               = errors.New("email inválido")
	ErrInvalidID                  = errors.New("ID inválido")
	ErrUserIDInvalid              = errors.New("ID de usuário inválido")
	ErrClientIDInvalid            = errors.New("ID de cliente inválido")
	ErrSupplierIDInvalid          = errors.New("ID de fornecedor inválido")
	ErrContactNotFound            = errors.New("contato não encontrado")
	ErrCreateContact              = errors.New("erro ao criar contato")
	ErrListUserContacts           = errors.New("erro ao listar contatos do usuário")
	ErrListClientContacts         = errors.New("erro ao listar contatos do cliente")
	ErrListSupplierContacts       = errors.New("erro ao listar contatos do fornecedor")
	ErrUpdateContact              = errors.New("erro ao atualizar contato")
	ErrDeleteContact              = errors.New("erro ao deletar contato")
	ErrCheckContact               = errors.New("erro ao verificar contato")
	ErrUpdateFailed               = errors.New("erro ao atualizar contato")
	ErrCheckBeforeUpdate          = errors.New("erro ao verificar existência do contato antes da atualização")
)
