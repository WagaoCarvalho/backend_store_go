package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/client_cpf/client"
)

type clientCpfHandler struct {
	service service.Client
	logger  *logger.LogAdapter
}

func NewClientCpfHandler(service service.Client, logger *logger.LogAdapter) *clientCpfHandler {
	return &clientCpfHandler{
		service: service,
		logger:  logger,
	}
}
