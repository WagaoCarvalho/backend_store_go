package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/category_relation"
)

type supplierCategoryRelationHandler struct {
	service service.SupplierCategoryRelation
	logger  *logger.LogAdapter
}

func NewSupplierCategoryRelationHandler(service service.SupplierCategoryRelation, logger *logger.LogAdapter) *supplierCategoryRelationHandler {
	return &supplierCategoryRelationHandler{
		service: service,
		logger:  logger,
	}
}
