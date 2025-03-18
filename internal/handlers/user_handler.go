package handlers

import (
	"fmt"
	"net/http"

	"github.com/WagaoCarvalho/backend_store_go/internal/services"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		utils.ErrorResponse(w, fmt.Errorf("erro ao buscar usuários: %w", err), http.StatusInternalServerError)
		return
	}

	response := utils.DefaultResponse{
		Data:   users,
		Status: http.StatusOK,
	}

	utils.ToJson(w, response)
}
