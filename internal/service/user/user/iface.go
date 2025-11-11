package services

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/user"

type User interface {
	iface.UserReader
	iface.UserWriter
	iface.UserStatus
}
