package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/category_relation"

type productCategoryRelationService struct {
	repo repo.ProductCategoryRelationRepo
}

func NewProductCategoryRelation(repo repo.ProductCategoryRelationRepo) ProductCategoryRelation {
	return &productCategoryRelationService{
		repo: repo,
	}
}
