package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/client_cpf/filter"
)

type clientCpfFilterHandler struct {
	service service.ClientCpfFilter
	logger  *logger.LogAdapter
}

func NewClientCpfFilterHandler(service service.ClientCpfFilter, logger *logger.LogAdapter) *clientCpfFilterHandler {
	return &clientCpfFilterHandler{
		service: service,
		logger:  logger,
	}
}
