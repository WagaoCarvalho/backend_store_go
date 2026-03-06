package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/address"

type Address interface {
	iface.AddressReader
	iface.AddressWriter
	iface.AddressStatus
}
