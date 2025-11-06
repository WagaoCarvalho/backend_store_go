package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/user"

type User interface {
	iface.UserWriter
	iface.UserReader
	iface.UserStatus
}
