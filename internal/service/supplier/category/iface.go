package services

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/supplier"

type SupplierCategory interface {
	iface.SupplierCategoryReader
	iface.SupplierCategoryWriter
}
