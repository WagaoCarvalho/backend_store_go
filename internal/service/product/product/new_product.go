package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/product"

type product struct {
	repo repo.Product
}

func NewProduct(repo repo.Product) ProductService {
	return &product{
		repo: repo,
	}
}
