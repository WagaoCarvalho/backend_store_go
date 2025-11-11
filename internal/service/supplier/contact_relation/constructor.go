package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/contact_relation"

type supplierContactRelationService struct {
	relationRepo repo.SupplierContactRelation
}

func NewSupplierContactRelation(repo repo.SupplierContactRelation) SupplierContactRelation {
	return &supplierContactRelationService{
		relationRepo: repo,
	}
}
