package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	iface "github.com/WagaoCarvalho/backend_store_go/internal/service/product/category"
)

type ProductCategory struct {
	service iface.ProductCategory
	logger  *logger.LogAdapter
}

func NewProductCategory(service iface.ProductCategory, logger *logger.LogAdapter) *ProductCategory {
	return &ProductCategory{
		service: service,
		logger:  logger,
	}
}
