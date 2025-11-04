package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user"
)

type User struct {
	service service.User
	logger  *logger.LogAdapter
}

func NewUser(service service.User, logger *logger.LogAdapter) *User {
	return &User{
		service: service,
		logger:  logger,
	}
}
