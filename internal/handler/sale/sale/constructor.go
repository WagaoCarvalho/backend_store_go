package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/sale/sale"
)

type SaleHandler struct {
	service service.SaleService
	logger  *logger.LogAdapter
}

func NewSaleHandler(service service.SaleService, logger *logger.LogAdapter) *SaleHandler {
	return &SaleHandler{
		service: service,
		logger:  logger,
	}
}
