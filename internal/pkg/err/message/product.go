package err

import (
	"errors"
)

var (
	ErrProductEnableDiscount     = errors.New("erro ao ativar desconto")
	ErrProductDisableDiscount    = errors.New("erro ao desativar desconto")
	ErrProductApplyDiscount      = errors.New("erro ao aplicar desconto")
	ErrProductDiscountNotAllowed = errors.New("erro desconto não permitido")
)
