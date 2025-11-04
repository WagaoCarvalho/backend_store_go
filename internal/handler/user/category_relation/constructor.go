package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/category_relation"
)

type UserCategoryRelation struct {
	service service.UserCategoryRelation
	logger  *logger.LogAdapter
}

func NewUserCategoryRelation(service service.UserCategoryRelation, logger *logger.LogAdapter) *UserCategoryRelation {
	return &UserCategoryRelation{
		service: service,
		logger:  logger,
	}
}
