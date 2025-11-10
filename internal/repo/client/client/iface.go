package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/client"

type Client interface {
	iface.ClientReader
	iface.ClientWriter
	iface.ClientStatus
	iface.ClientFilter
}
