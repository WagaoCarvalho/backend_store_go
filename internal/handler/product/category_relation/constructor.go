package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	iface "github.com/WagaoCarvalho/backend_store_go/internal/service/product/category_relation"
)

type ProductCategoryRelation struct {
	productCategoryRelation iface.ProductCategoryRelation
	logger                  *logger.LogAdapter
}

func NewProductCategoryRelation(productCategoryRelation iface.ProductCategoryRelation, logger *logger.LogAdapter) *ProductCategoryRelation {
	return &ProductCategoryRelation{
		productCategoryRelation: productCategoryRelation,
		logger:                  logger,
	}
}
