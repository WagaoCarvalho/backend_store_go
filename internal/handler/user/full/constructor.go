package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/full"
)

type userHandler struct {
	service service.UserFull
	logger  *logger.LogAdapter
}

func NewUserFullHandler(service service.UserFull, logger *logger.LogAdapter) *userHandler {
	return &userHandler{
		service: service,
		logger:  logger,
	}
}
