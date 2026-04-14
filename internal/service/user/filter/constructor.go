package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/filter"

type userFilterService struct {
	repo repo.UserFilter
}

func NewUserFilterService(repo repo.UserFilter) UserFilter {
	return &userFilterService{
		repo: repo,
	}
}
