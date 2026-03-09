package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/address/filter"
)

type addressFilterHandler struct {
	service service.AddressFilter
	logger  *logger.LogAdapter
}

func NewAddressFilterHandler(service service.AddressFilter, logger *logger.LogAdapter) *addressFilterHandler {
	return &addressFilterHandler{
		service: service,
		logger:  logger,
	}
}
