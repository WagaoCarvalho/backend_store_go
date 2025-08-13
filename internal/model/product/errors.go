package models

import "errors"

var (
	ErrInvalidProductName        = errors.New("nome do produto é obrigatório")
	ErrInvalidManufacturer       = errors.New("fabricante é obrigatório")
	ErrInvalidCostPrice          = errors.New("preço de custo deve ser maior ou igual a zero")
	ErrInvalidSalePrice          = errors.New("preço de venda deve ser maior ou igual a zero")
	ErrNegativeStock             = errors.New("quantidade em estoque não pode ser negativa")
	ErrInvalidBarcode            = errors.New("código de barras inválido")
	ErrSupplierRequired          = errors.New("ID do fornecedor deve ser informado")
	ErrInactiveProductNotAllowed = errors.New("não permitido, produto inativo")
)
