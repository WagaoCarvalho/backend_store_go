package services

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/supplier"

type Supplier interface {
	iface.SupplierReader
	iface.SupplierWriter
	iface.SupplierStatus
}
