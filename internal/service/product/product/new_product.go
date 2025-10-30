package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/product"

type product struct {
	repo repo.ProductRepo
}

func NewProduct(repo repo.ProductRepo) ProductService {
	return &product{
		repo: repo,
	}
}
