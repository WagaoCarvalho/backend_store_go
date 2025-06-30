package services

import "errors"

var (
	ErrInvalidEmail   = errors.New("email inválido")
	ErrCreateUser     = errors.New("erro ao criar usuário")
	ErrGetUser        = errors.New("erro ao buscar usuário")
	ErrGetVersion     = errors.New("erro ao buscar usuário")
	ErrUpdateUser     = errors.New("erro ao atualizar usuário")
	ErrDeleteUser     = errors.New("erro ao deletar usuário")
	ErrInvalidVersion = errors.New("versão inválida")
)
