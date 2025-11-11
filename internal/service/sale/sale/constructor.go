package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/sale/sale"

type saleService struct {
	repo repo.SaleRepo
}

func NewSaleService(repo repo.SaleRepo) SaleService {
	return &saleService{
		repo: repo,
	}
}
