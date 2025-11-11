package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/category"

type supplierCategoryService struct {
	repo repo.SupplierCategory
}

func NewSupplierCategory(repo repo.SupplierCategory) SupplierCategory {
	return &supplierCategoryService{
		repo: repo,
	}
}
