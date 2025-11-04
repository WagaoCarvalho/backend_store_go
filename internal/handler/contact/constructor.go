package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/contact"
)

type Contact struct {
	service service.Contact
	logger  *logger.LogAdapter
}

func NewContact(service service.Contact, logger *logger.LogAdapter) *Contact {
	return &Contact{
		service: service,
		logger:  logger,
	}
}
