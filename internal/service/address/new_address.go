package services

import (
	repoAddress "github.com/WagaoCarvalho/backend_store_go/internal/repo/address"
	repoClient "github.com/WagaoCarvalho/backend_store_go/internal/repo/client/client"
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
	repoUser "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
)

type address struct {
	repoAddress  repoAddress.Address
	repoClient   repoClient.Client
	repoUser     repoUser.User
	repoSupplier repoSupplier.Supplier
}

func NewAddress(
	repoAddress repoAddress.Address,
	repoClient repoClient.Client,
	repoUser repoUser.User,
	repoSupplier repoSupplier.Supplier,
) Address {
	return &address{
		repoAddress:  repoAddress,
		repoClient:   repoClient,
		repoUser:     repoUser,
		repoSupplier: repoSupplier,
	}
}
