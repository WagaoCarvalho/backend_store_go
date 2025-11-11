package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/full"
)

type supplierHandler struct {
	service service.SupplierFullService
	logger  *logger.LogAdapter
}

func NewSupplierFull(service service.SupplierFullService, logger *logger.LogAdapter) *supplierHandler {
	return &supplierHandler{
		service: service,
		logger:  logger,
	}
}
