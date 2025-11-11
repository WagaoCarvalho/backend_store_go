package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/category_relation"

type userCategoryRelationService struct {
	relationRepo repo.UserCategoryRelation
}

func NewUserCategoryRelationService(repo repo.UserCategoryRelation) UserCategoryRelation {
	return &userCategoryRelationService{
		relationRepo: repo,
	}
}
