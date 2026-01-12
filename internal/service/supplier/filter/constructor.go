package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/filter"

type supplierFilterService struct {
	repo repo.SupplierFilter
}

func NewSupplierFilterService(repo repo.SupplierFilter) SupplierFilter {
	return &supplierFilterService{
		repo: repo,
	}
}
