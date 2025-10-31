package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/sale/item"
)

type SaleItemHandler struct {
	service service.SaleItemService
	logger  *logger.LogAdapter
}

func NewSaleItemHandler(service service.SaleItemService, logger *logger.LogAdapter) *SaleItemHandler {
	return &SaleItemHandler{
		service: service,
		logger:  logger,
	}
}
