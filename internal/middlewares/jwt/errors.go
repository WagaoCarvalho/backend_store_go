package middlewares

import "errors"

var (
	ErrTokenMissing         = errors.New("token ausente")
	ErrTokenInvalidFormat   = errors.New("formato de token inválido")
	ErrTokenRevoked         = errors.New("token revogado")
	ErrTokenInvalid         = errors.New("token inválido")
	ErrTokenExpired         = errors.New("token expirado")
	ErrInvalidSignature     = errors.New("assinatura inválida")
	ErrInvalidSigningMethod = errors.New("método de assinatura inválido")
	ErrInvalidExpClaim      = errors.New("campo exp ausente ou inválido")
	ErrInternalAuth         = errors.New("erro interno de autenticação")
)
