package handler

import (
	iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/product"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
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
