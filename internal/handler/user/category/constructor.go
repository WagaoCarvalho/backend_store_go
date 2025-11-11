package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/category"
)

type userCategoryHandler struct {
	service service.UserCategory
	logger  *logger.LogAdapter
}

func NewUserCategoryHandler(service service.UserCategory, logger *logger.LogAdapter) *userCategoryHandler {
	return &userCategoryHandler{
		service: service,
		logger:  logger,
	}
}
