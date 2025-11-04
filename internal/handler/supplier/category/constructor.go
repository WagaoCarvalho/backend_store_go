package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/category"
)

type SupplierCategory struct {
	service service.SupplierCategory
	logger  *logger.LogAdapter
}

func NewSupplierCategory(service service.SupplierCategory, logger *logger.LogAdapter) *SupplierCategory {
	return &SupplierCategory{
		service: service,
		logger:  logger,
	}
}
