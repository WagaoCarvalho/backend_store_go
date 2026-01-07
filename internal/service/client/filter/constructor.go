package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/filter"

type clientFiltertService struct {
	repo repo.ClientFilter
}

func NewClientFilterService(repo repo.ClientFilter) ClientFilter {
	return &clientFiltertService{
		repo: repo,
	}
}
