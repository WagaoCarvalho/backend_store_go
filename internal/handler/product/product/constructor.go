package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product/product"
)

type Product struct {
	service service.ProductService
	logger  *logger.LogAdapter
}

func NewProduct(service service.ProductService, logger *logger.LogAdapter) *Product {
	return &Product{
		service: service,
		logger:  logger,
	}
}
