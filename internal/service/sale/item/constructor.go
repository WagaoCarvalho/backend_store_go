package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/sale/item"

type saleItem struct {
	repo repo.SaleItemRepo
}

func NewItemSale(repo repo.SaleItemRepo) SaleItemService {
	return &saleItem{
		repo: repo,
	}
}
