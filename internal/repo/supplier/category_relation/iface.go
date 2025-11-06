package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/supplier"

type SupplierCategoryRelation interface {
	iface.SupplierCategoryRelationReader
	iface.SupplierCategoryRelationWriter
}
