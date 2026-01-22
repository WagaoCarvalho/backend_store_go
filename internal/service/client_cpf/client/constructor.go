package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/client_cpf/client"

type clientCPfService struct {
	repo repo.ClientCpf
}

func NewClientCpfService(repo repo.ClientCpf) Client {
	return &clientCPfService{
		repo: repo,
	}
}
