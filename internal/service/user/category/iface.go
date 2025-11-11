package services

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/user"

type UserCategory interface {
	iface.UserCategoryReader
	iface.UserCategoryWriter
}
