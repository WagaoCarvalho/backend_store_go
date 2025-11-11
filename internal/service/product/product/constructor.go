package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/product"

type productService struct {
	repo repo.Product
}

func NewProductService(repo repo.Product) ProductService {
	return &productService{
		repo: repo,
	}
}
