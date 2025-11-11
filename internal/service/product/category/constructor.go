package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/category"

type productCategoryService struct {
	repo repo.ProductCategory
}

func NewProductCategoryService(repo repo.ProductCategory) ProductCategory {
	return &productCategoryService{
		repo: repo,
	}
}
