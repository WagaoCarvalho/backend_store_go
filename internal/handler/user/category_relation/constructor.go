package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/category_relation"
)

type userCategoryRelationHandler struct {
	service service.UserCategoryRelation
	logger  *logger.LogAdapter
}

func NewUserCategoryRelationHandler(service service.UserCategoryRelation, logger *logger.LogAdapter) *userCategoryRelationHandler {
	return &userCategoryRelationHandler{
		service: service,
		logger:  logger,
	}
}
