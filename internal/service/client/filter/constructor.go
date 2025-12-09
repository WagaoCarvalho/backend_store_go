package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/filter"

type clienFiltertService struct {
	repo repo.ClientFilter
}

func NewClientFilterService(repo repo.ClientFilter) ClientFilter {
	return &clienFiltertService{
		repo: repo,
	}
}
