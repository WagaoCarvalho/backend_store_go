package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/filter"

type productFilterService struct {
	repo repo.ProductFilter
}

func NewProductFilterService(repo repo.ProductFilter) ProductFilter {
	return &productFilterService{
		repo: repo,
	}
}
