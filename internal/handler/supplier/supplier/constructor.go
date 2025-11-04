package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier"
)

type SupplierHandler struct {
	service service.Supplier
	logger  *logger.LogAdapter
}

func NewSupplierHandler(service service.Supplier, logger *logger.LogAdapter) *SupplierHandler {
	return &SupplierHandler{
		service: service,
		logger:  logger,
	}
}
