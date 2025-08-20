package services

import "errors"

var (
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
	ErrEnableProduct             = errors.New("erro ao ativar produto")
	ErrDisableProduct            = errors.New("erro ao desativar produto")
	ErrUpdateStock               = errors.New("falha ao atualizar estoque do produto")
	ErrGetStock                  = errors.New("erro ao obter o estoque do produto")

	ErrEnableDiscount  = errors.New("erro ao ativar desconto")
	ErrDisableDiscount = errors.New("erro ao desativar desconto")
	ErrApplyDiscount   = errors.New("erro ao aplicar desconto")
)
