package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	iface "github.com/WagaoCarvalho/backend_store_go/internal/service/product/category"
)

type productCategoryHandler struct {
	service iface.ProductCategory
	logger  *logger.LogAdapter
}

func NewProductCategoryHandler(service iface.ProductCategory, logger *logger.LogAdapter) *productCategoryHandler {
	return &productCategoryHandler{
		service: service,
		logger:  logger,
	}
}
