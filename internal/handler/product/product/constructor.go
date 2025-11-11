package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product/product"
)

type productHandler struct {
	service service.ProductService
	logger  *logger.LogAdapter
}

func NewProductHandler(service service.ProductService, logger *logger.LogAdapter) *productHandler {
	return &productHandler{
		service: service,
		logger:  logger,
	}
}
