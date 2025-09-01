package handler

import (
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/login"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/login"
)

type LoginHandler struct {
	service service.LoginService
	logger  *logger.LogAdapter
}

func NewLoginHandler(service service.LoginService, logger *logger.LogAdapter) *LoginHandler {
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

	// ler DTO
	var credentialsDTO dto.LoginCredentialsDTO
	if err := utils.FromJSON(r.Body, &credentialsDTO); err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogParseJSONError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, fmt.Errorf("dados inválidos"), http.StatusBadRequest)
		return
	}

	// converter DTO para model
	credentials := credentialsDTO.ToModel()

	// chamar service
	authResp, err := h.service.Login(r.Context(), *credentials)
	if err != nil {
		h.logger.Warn(r.Context(), ref+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	// converter model para DTO de resposta
	authRespDTO := dto.ToAuthResponseDTO(authResp)

	h.logger.Info(r.Context(), ref+logger.LogLoginSuccess, nil)

	// retornar JSON
	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Login realizado com sucesso",
		Data:    authRespDTO, // compatível com JSON
	})
}
