package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/address"
)

type addressHandler struct {
	service service.Address
	logger  *logger.LogAdapter
}

func NewAddressHandler(service service.Address, logger *logger.LogAdapter) *addressHandler {
	return &addressHandler{
		service: service,
		logger:  logger,
	}
}
