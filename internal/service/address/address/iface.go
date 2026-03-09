package services

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/address"

type Address interface {
	iface.AddressReader
	iface.AddressStatus
	iface.AddressWriter
}
