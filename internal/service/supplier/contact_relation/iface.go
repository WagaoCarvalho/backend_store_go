package services

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/supplier"

type SupplierContactRelation interface {
	iface.SupplierContactRelationReader
	iface.SupplierContactRelationWriter
}
