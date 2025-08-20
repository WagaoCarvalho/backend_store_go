package models

import "errors"

var (
	ErrInvalidProductName        = errors.New("produto inválido")
	ErrInvalidManufacturer       = errors.New("fabricante inválido")
	ErrInvalidCostPrice          = errors.New("preço de custo não pode ser negativo")
	ErrInvalidSalePrice          = errors.New("preço de venda não pode ser negativo")
	ErrSalePriceBelowCost        = errors.New("preço de venda não pode ser menor que o preço de custo")
	ErrNegativeStock             = errors.New("estoque não pode ser negativo")
	ErrInvalidBarcode            = errors.New("código de barras inválido")
	ErrSupplierRequired          = errors.New("fornecedor é obrigatório")
	ErrInactiveProductNotAllowed = errors.New("produto inativo não é permitido")
	ErrNegativeDiscount          = errors.New("percentual de desconto não pode ser negativo")
	ErrDiscountAboveLimit        = errors.New("percentual de desconto não pode exceder 100%")
	ErrInvalidDiscountRange      = errors.New("percentual mínimo de desconto não pode ser maior que o percentual máximo")
)
