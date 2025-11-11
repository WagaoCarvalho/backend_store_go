package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/sale/sale"
)

type saleHandler struct {
	service service.SaleService
	logger  *logger.LogAdapter
}

func NewSaleHandler(service service.SaleService, logger *logger.LogAdapter) *saleHandler {
	return &saleHandler{
		service: service,
		logger:  logger,
	}
}
