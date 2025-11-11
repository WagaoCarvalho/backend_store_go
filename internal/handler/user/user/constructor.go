package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user"
)

type userHandler struct {
	service service.User
	logger  *logger.LogAdapter
}

func NewUserHandler(service service.User, logger *logger.LogAdapter) *userHandler {
	return &userHandler{
		service: service,
		logger:  logger,
	}
}
