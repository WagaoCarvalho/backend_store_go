package err

import (
	"errors"
)

var (
	ErrCreateProduct       = errors.New("erro ao criar produto")
	ErrGetProduct          = errors.New("erro ao buscar produto")
	ErrGetProducts         = errors.New("erro ao buscar produtos")
	ErrUpdateProduct       = errors.New("erro ao atualizar produto")
	ErrDeleteProduct       = errors.New("erro ao excluir produto")
	ErrProductNotFound     = errors.New("produto não encontrado")
	ErrFetchProductVersion = errors.New("erro ao buscar versão do produto")
	ErrDisableProduct      = errors.New("erro ao desabilitar produto")
	ErrEnableProduct       = errors.New("erro ao ativar produto")
	ErrVersionConflict     = errors.New("conflito de versão")
	ErrUpdateStock         = errors.New("falha ao atualizar estoque do produto")

	ErrEnableDiscount  = errors.New("erro ao ativar desconto")
	ErrDisableDiscount = errors.New("erro ao desativar desconto")
	ErrApplyDiscount   = errors.New("erro ao aplicar desconto")

	ErrDiscountNotAllowed = errors.New("erro desconto não permitido")
)
