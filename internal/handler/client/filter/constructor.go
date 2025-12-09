package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/client/filter"
)

type clientFilterHandler struct {
	service service.ClientFilter
	logger  *logger.LogAdapter
}

func NewClientFilterHandler(service service.ClientFilter, logger *logger.LogAdapter) *clientFilterHandler {
	return &clientFilterHandler{
		service: service,
		logger:  logger,
	}
}
