package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/category"
)

type supplierCategoryHandler struct {
	service service.SupplierCategory
	logger  *logger.LogAdapter
}

func NewSupplierCategoryHandler(service service.SupplierCategory, logger *logger.LogAdapter) *supplierCategoryHandler {
	return &supplierCategoryHandler{
		service: service,
		logger:  logger,
	}
}
