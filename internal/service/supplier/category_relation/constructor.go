package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/category_relation"

type supplierCategoryRelationService struct {
	relationRepo repo.SupplierCategoryRelation
}

func NewSupplierCategoryRelationService(repository repo.SupplierCategoryRelation) SupplierCategoryRelation {
	return &supplierCategoryRelationService{relationRepo: repository}
}
