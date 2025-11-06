package services

import (
	repoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
	repoUser "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/client/client"
)

type address struct {
	address  Address
	client   service.Client
	user     repoUser.User
	supplier repoSupplier.Supplier
}

func NewAddress(
	repoAddress Address,
	repoClient service.Client,
	repoUser repoUser.User,
	repoSupplier repoSupplier.Supplier,
) Address {
	return &address{
		address:  repoAddress,
		client:   repoClient,
		user:     repoUser,
		supplier: repoSupplier,
	}
}
