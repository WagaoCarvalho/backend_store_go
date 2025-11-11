package services

import (
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoClient "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/client"
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
	repoUser "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
)

type addressService struct {
	addressRepo  repoAddress.Address
	clientRepo   repoClient.Client
	userRepo     repoUser.User
	supplierRepo repoSupplier.Supplier
}

// Retorna a interface da camada de serviço
func NewAddress(
	addressRepo repoAddress.Address,
	clientRepo repoClient.Client,
	userRepo repoUser.User,
	supplierRepo repoSupplier.Supplier,
) Address {
	return &addressService{
		addressRepo:  addressRepo,
		clientRepo:   clientRepo,
		userRepo:     userRepo,
		supplierRepo: supplierRepo,
	}
}
