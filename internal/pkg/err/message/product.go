package err

import (
	"errors"
)

var (
	ErrCreateProduct          = errors.New("erro ao criar produto")
	ErrGetProduct             = errors.New("erro ao buscar produto")
	ErrGetProducts            = errors.New("erro ao buscar produtos")
	ErrUpdateProduct          = errors.New("erro ao atualizar produto")
	ErrDeleteProduct          = errors.New("erro ao excluir produto")
	ErrProductNotFound        = errors.New("produto não encontrado")
	ErrFetchProductVersion    = errors.New("erro ao buscar versão do produto")
	ErrDisableProduct         = errors.New("erro ao desabilitar produto")
	ErrEnableProduct          = errors.New("erro ao ativar produto")
	ErrProductVersionConflict = errors.New("conflito de versão")
	ErrUpdateStock            = errors.New("falha ao atualizar estoque do produto")

	ErrEnableDiscount  = errors.New("erro ao ativar desconto")
	ErrDisableDiscount = errors.New("erro ao desativar desconto")
	ErrApplyDiscount   = errors.New("erro ao aplicar desconto")

	ErrDiscountNotAllowed = errors.New("erro desconto não permitido")

	ErrProductFetch               = errors.New("erro ao obter produtos")
	ErrProductFetchByID           = errors.New("erro ao obter produto")
	ErrProductFetchByName         = errors.New("erro ao obter produtos por nome")
	ErrProductFetchByManufacturer = errors.New("erro ao obter produtos por fabricante")
	ErrProductCreateNameRequired  = errors.New("validação falhou: nome do produto é obrigatório")
	ErrProductCreateCostPrice     = errors.New("validação falhou: preço de custo deve ser positivo")

	ErrInvalidProduct            = errors.New("produto inválido")
	ErrProductCreateManufacturer = errors.New("validação falhou: fabricante é obrigatório")
	ErrProductCreatePriceLogic   = errors.New("validação falhou: preço de venda deve ser maior que o preço de custo")
	ErrProductUpdate             = errors.New("erro ao atualizar produto")
	ErrProductDelete             = errors.New("erro ao deletar produto")
	ErrProductFetchByCostPrice   = errors.New("erro ao obter produtos por faixa de preço de custo")
	ErrProductFetchBySalePrice   = errors.New("erro ao obter produtos por faixa de preço de venda")
	ErrProductLowStock           = errors.New("erro ao buscar produtos com estoque baixo")
	ErrInvalidVersion            = errors.New("versão inválida")
	ErrGetStock                  = errors.New("erro ao obter o estoque do produto")
)
