package handler

import (
	iface "github.com/WagaoCarvalho/backend_store_go/internal/iface/product"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
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
