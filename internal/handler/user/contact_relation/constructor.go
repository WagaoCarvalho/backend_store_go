package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/contact_relation"
)

type userContactRelationHandler struct {
	service service.UserContactRelation
	logger  *logger.LogAdapter
}

func NewUserContactRelationHandler(service service.UserContactRelation, logger *logger.LogAdapter) *userContactRelationHandler {
	return &userContactRelationHandler{
		service: service,
		logger:  logger,
	}
}
