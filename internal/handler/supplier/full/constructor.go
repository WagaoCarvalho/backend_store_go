package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/full"
)

type Supplier struct {
	service service.SupplierFull
	logger  *logger.LogAdapter
}

func NewSupplierFull(service service.SupplierFull, logger *logger.LogAdapter) *Supplier {
	return &Supplier{
		service: service,
		logger:  logger,
	}
}
