package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/client_cpf/filter"

type clientCpfFiltertService struct {
	repo repo.ClientCpfFilter
}

func NewClientCpfFilterService(repo repo.ClientCpfFilter) ClientCpfFilter {
	return &clientCpfFiltertService{
		repo: repo,
	}
}
