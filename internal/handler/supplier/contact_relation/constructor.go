package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/contact_relation"
)

type supplierContactRelationHandler struct {
	service service.SupplierContactRelation
	logger  *logger.LogAdapter
}

func NewSupplierContactRelationHandler(service service.SupplierContactRelation, logger *logger.LogAdapter) *supplierContactRelationHandler {
	return &supplierContactRelationHandler{
		service: service,
		logger:  logger,
	}
}
