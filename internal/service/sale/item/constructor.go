package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/sale/item"

type saleItemService struct {
	repo repo.SaleItemRepo
}

func NewItemSaleService(repo repo.SaleItemRepo) SaleItemService {
	return &saleItemService{
		repo: repo,
	}
}
