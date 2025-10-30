package services

import "github.com/WagaoCarvalho/backend_store_go/internal/iface"

type ProductService interface {
	iface.ProductReader
	iface.ProductWriter
	iface.ProductStock
	iface.ProductDiscount
	iface.ProductStatus
}
