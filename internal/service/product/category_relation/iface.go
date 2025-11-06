package services

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/product"

type ProductCategoryRelation interface {
	iface.ProductCategoryRelationReader
	iface.ProductCategoryRelationWriter
}
