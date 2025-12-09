package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/filter"

type clienFiltertService struct {
	repo repo.Client
}

func NewClientFilterService(repo repo.Client) Client {
	return &clienFiltertService{
		repo: repo,
	}
}
