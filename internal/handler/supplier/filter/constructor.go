package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/filter"
)

type supplierFilterHandler struct {
	service service.SupplierFilter
	logger  *logger.LogAdapter
}

func NewSupplierFilterHandler(service service.SupplierFilter, logger *logger.LogAdapter) *supplierFilterHandler {
	return &supplierFilterHandler{
		service: service,
		logger:  logger,
	}
}
