package services

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/product"

type ProductCategory interface {
	iface.ProductCategoryReader
	iface.ProductCategoryWriter
}
