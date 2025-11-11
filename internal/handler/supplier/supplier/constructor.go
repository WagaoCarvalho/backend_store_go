package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier"
)

type supplierHandler struct {
	service service.Supplier
	logger  *logger.LogAdapter
}

func NewSupplierHandler(service service.Supplier, logger *logger.LogAdapter) *supplierHandler {
	return &supplierHandler{
		service: service,
		logger:  logger,
	}
}
