package repo

import iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/product"

type ProductCategoryRelationRepo interface {
	iface.ProductCategoryRelationReader
	iface.ProductCategoryRelationWriter
	iface.ProductCategoryRelationChecker
}
