package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/sale/filter"

type saleFilterService struct {
	repo repo.SaleFilter
}

func NewSaleFilterService(repo repo.SaleFilter) SaleFilter {
	return &saleFilterService{
		repo: repo,
	}
}
