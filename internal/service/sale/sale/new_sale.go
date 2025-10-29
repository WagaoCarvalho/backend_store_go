package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/sale"

type sale struct {
	repo repo.SaleRepo
}

func NewSale(repo repo.SaleRepo) SaleService {
	return &sale{
		repo: repo,
	}
}
