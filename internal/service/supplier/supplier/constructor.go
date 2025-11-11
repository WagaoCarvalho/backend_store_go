package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"

type supplierService struct {
	repo repo.Supplier
}

func NewSupplierService(repo repo.Supplier) Supplier {
	return &supplierService{
		repo: repo,
	}
}
