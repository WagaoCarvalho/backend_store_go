package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/product/filter"
)

type productFilterHandler struct {
	service service.ProductFilter
	logger  *logger.LogAdapter
}

func NewProductFilterHandler(service service.ProductFilter, logger *logger.LogAdapter) *productFilterHandler {
	return &productFilterHandler{
		service: service,
		logger:  logger,
	}
}
