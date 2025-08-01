package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	auth "github.com/WagaoCarvalho/backend_store_go/internal/auth/login"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/login"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
)

type LoginHandler struct {
	service auth.LoginService
	logger  *logger.LoggerAdapter
}

func NewLoginHandler(service auth.LoginService, logger *logger.LoggerAdapter) *LoginHandler {
	return &LoginHandler{
		service: service,
		logger:  logger,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	const ref = "[LoginHandler - Login] "

	h.logger.Info(r.Context(), ref+logger.LogLoginInit, nil)

	if r.Method != http.MethodPost {
		h.logger.Warn(r.Context(), ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	var credentials models.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), credentials)
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	h.logger.Info(r.Context(), ref+logger.LogLoginSuccess, nil)

	utils.ToJson(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Login realizado com sucesso",
		Data:    map[string]string{"token": token},
	})
}
