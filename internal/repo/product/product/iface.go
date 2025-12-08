package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/product"

type Product interface {
	iface.ProductReader
	iface.ProductWriter
	iface.ProductFilter
	iface.ProductStock
	iface.ProductDiscount
	iface.ProductStatus
	iface.ProductGetAll
	iface.ProductChecker
	iface.ProductVersion
}
