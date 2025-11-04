package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/contact_relation"
)

type UserContactRelation struct {
	service service.UserContactRelation
	logger  *logger.LogAdapter
}

func NewUserContactRelation(service service.UserContactRelation, logger *logger.LogAdapter) *UserContactRelation {
	return &UserContactRelation{
		service: service,
		logger:  logger,
	}
}
