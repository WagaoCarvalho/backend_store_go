package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/category"
)

type UserCategory struct {
	service service.UserCategory
	logger  *logger.LogAdapter
}

func NewUserCategory(service service.UserCategory, logger *logger.LogAdapter) *UserCategory {
	return &UserCategory{
		service: service,
		logger:  logger,
	}
}
