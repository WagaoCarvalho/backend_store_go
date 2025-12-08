package services

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/product"

type ProductService interface {
	iface.ProductReader
	iface.ProductWriter
	iface.ProductStock
	iface.ProductDiscount
	iface.ProductStatus
	iface.ProductFilter
	iface.ProductVersion
}
