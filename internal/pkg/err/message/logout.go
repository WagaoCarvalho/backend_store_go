package err

import "errors"

var (
	ErrInvalidSigningMethod = errors.New("método de assinatura inválido")
	ErrTokenValidation      = errors.New("erro ao validar token")
	ErrInvalidToken         = errors.New("token inválido")
	ErrClaimConversion      = errors.New("não foi possível obter claims")
	ErrClaimExpInvalid      = errors.New("claim 'exp' ausente ou inválida")
	ErrTokenExpired         = errors.New("token já expirado")
	ErrBlacklistAdd         = errors.New("erro ao realizar logout")
)
