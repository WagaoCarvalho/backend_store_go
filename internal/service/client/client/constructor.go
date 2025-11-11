package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/client"

type clientService struct {
	repo repo.Client
}

func NewClientService(repo repo.Client) Client {
	return &clientService{
		repo: repo,
	}
}
