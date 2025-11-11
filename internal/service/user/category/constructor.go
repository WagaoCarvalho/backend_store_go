package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/category"

type userCategoryService struct {
	repo repo.UserCategory
}

func NewUserCategoryService(repo repo.UserCategory) UserCategory {
	return &userCategoryService{
		repo: repo,
	}
}
