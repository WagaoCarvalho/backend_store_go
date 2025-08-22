package err

import "errors"

var (
	ErrInvalidEmailFormat = errors.New("formato de email inválido")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUserFetchFailed    = errors.New("erro ao buscar usuário")
	ErrTokenGeneration    = errors.New("erro ao gerar token de acesso")
	ErrAccountDisabled    = errors.New("conta desativada")
)
