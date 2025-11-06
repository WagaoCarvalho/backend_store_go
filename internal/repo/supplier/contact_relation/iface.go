package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/supplier"

type SupplierContactRelation interface {
	iface.SupplierContactRelationWriter
	iface.SupplierContactRelationReader
}
