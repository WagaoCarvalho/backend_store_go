package services

import repoClient "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/client"

type clientService struct {
	repo repoClient.Client
}

func NewClientService(repo Client) Client {
	return &clientService{
		repo: repo,
	}
}
