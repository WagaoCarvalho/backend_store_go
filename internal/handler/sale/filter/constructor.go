package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/sale/filter"
)

type saleFilterHandler struct {
	service service.SaleFilter
	logger  *logger.LogAdapter
}

func NewSaleFilterHandler(service service.SaleFilter, logger *logger.LogAdapter) *saleFilterHandler {
	return &saleFilterHandler{
		service: service,
		logger:  logger,
	}
}
