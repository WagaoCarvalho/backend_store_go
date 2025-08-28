package err

import "errors"

var (
	ErrEmailFormat     = errors.New("formato de email inválido")
	ErrCredentials     = errors.New("credenciais inválidas")
	ErrTokenGeneration = errors.New("erro ao gerar token de acesso")
	ErrAccountDisabled = errors.New("conta desativada")
)
