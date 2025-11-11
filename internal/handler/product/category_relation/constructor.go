package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	iface "github.com/WagaoCarvalho/backend_store_go/internal/service/product/category_relation"
)

type productCategoryRelationHandler struct {
	productCategoryRelation iface.ProductCategoryRelation
	logger                  *logger.LogAdapter
}

func NewProductCategoryRelationHandler(
	productCategoryRelation iface.ProductCategoryRelation,
	logger *logger.LogAdapter,
) *productCategoryRelationHandler {
	return &productCategoryRelationHandler{
		productCategoryRelation: productCategoryRelation,
		logger:                  logger,
	}
}
