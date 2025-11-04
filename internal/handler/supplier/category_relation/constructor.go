package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/category_relation"
)

type SupplierCategoryRelation struct {
	service service.SupplierCategoryRelation
	logger  *logger.LogAdapter
}

func NewSupplierCategoryRelation(service service.SupplierCategoryRelation, logger *logger.LogAdapter) *SupplierCategoryRelation {
	return &SupplierCategoryRelation{
		service: service,
		logger:  logger,
	}
}
