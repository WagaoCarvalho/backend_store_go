package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/filter"

type productFiltertService struct {
	repo repo.ProductFilter
}

func NewProductFilterService(repo repo.ProductFilter) ProductFilter {
	return &productFiltertService{
		repo: repo,
	}
}
