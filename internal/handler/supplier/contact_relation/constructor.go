package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/contact_relation"
)

type SupplierContactRelation struct {
	service service.SupplierContactRelation
	logger  *logger.LogAdapter
}

func NewSupplierContactRelation(service service.SupplierContactRelation, logger *logger.LogAdapter) *SupplierContactRelation {
	return &SupplierContactRelation{
		service: service,
		logger:  logger,
	}
}
