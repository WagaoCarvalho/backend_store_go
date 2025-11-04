package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/full"
)

type UserHandler struct {
	service service.UserFull
	logger  *logger.LogAdapter
}

func NewUserFullHandler(service service.UserFull, logger *logger.LogAdapter) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}
