package services

import repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/address/filter"

type addressFiltertService struct {
	repo repo.AddressFilter
}

func NewAddressFilterService(repo repo.AddressFilter) AddressFilter {
	return &addressFiltertService{
		repo: repo,
	}
}
