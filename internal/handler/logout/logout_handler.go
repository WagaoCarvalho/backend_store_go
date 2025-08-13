package handler

import (
	"net/http"
	"strings"

	service "github.com/WagaoCarvalho/backend_store_go/internal/auth/logout"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
)

type LogoutHandler struct {
	service service.LogoutService
	logger  *logger.LoggerAdapter
}

func NewLogoutHandler(service service.LogoutService, logger *logger.LoggerAdapter) *LogoutHandler {
	return &LogoutHandler{
		service: service,
		logger:  logger,
	}
}

func (h *LogoutHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorResponse(w, nil, http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.ErrorResponse(w, nil, http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		utils.ErrorResponse(w, nil, http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]

	err := h.service.Logout(r.Context(), tokenString)
	if err != nil {
		h.logger.Error(r.Context(), err, "Erro ao invalidar token na blacklist", nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Logout realizado com sucesso",
	})
}
