package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/client/client"
)

type clientHandler struct {
	service service.Client
	logger  *logger.LogAdapter
}

func NewClientHandler(service service.Client, logger *logger.LogAdapter) *clientHandler {
	return &clientHandler{
		service: service,
		logger:  logger,
	}
}
