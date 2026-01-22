package services

import (
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoClientCpf "github.com/WagaoCarvalho/backend_store_go/internal/repo/client_cpf/client"
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
	repoUser "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
)

type addressService struct {
	addressRepo   repoAddress.Address
	clientCpfRepo repoClientCpf.ClientCpf
	userRepo      repoUser.User
	supplierRepo  repoSupplier.Supplier
}

// Retorna a interface da camada de serviço
func NewAddressService(
	addressRepo repoAddress.Address,
	clientCpfRepo repoClientCpf.ClientCpf,
	userRepo repoUser.User,
	supplierRepo repoSupplier.Supplier,
) Address {
	return &addressService{
		addressRepo:   addressRepo,
		clientCpfRepo: clientCpfRepo,
		userRepo:      userRepo,
		supplierRepo:  supplierRepo,
	}
}
