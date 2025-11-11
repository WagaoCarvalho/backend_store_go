package handler

import (
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/contact"
)

type contactHandler struct {
	service service.Contact
	logger  *logger.LogAdapter
}

func NewContactHandler(service service.Contact, logger *logger.LogAdapter) *contactHandler {
	return &contactHandler{
		service: service,
		logger:  logger,
	}
}
