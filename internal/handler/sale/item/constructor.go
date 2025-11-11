package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/sale/item"
)

type saleItemHandler struct {
	service service.SaleItemService
	logger  *logger.LogAdapter
}

func NewSaleItemHandler(service service.SaleItemService, logger *logger.LogAdapter) *saleItemHandler {
	return &saleItemHandler{
		service: service,
		logger:  logger,
	}
}
