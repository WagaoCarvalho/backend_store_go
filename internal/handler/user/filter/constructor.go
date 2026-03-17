package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/filter"
)

type userFilterHandler struct {
	service service.UserFilter
	logger  *logger.LogAdapter
}

func NewUserFilterHandler(service service.UserFilter, logger *logger.LogAdapter) *userFilterHandler {
	return &userFilterHandler{
		service: service,
		logger:  logger,
	}
}
