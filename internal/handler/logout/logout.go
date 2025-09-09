package handler

import (
	"net/http"
	"strings"

	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/logout"
)

type LogoutHandler struct {
	service service.LogoutService
	logger  *logger.LogAdapter
}

func NewLogoutHandler(service service.LogoutService, logger *logger.LogAdapter) *LogoutHandler {
	return &LogoutHandler{
		service: service,
		logger:  logger,
	}
}

func (h *LogoutHandler) Logout(w http.ResponseWriter, r *http.Request) {
	const ref = "[LogoutHandler - Logout] "
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+"Método não permitido", map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, nil, http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+"Iniciando logout", nil)

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.logger.Warn(ctx, ref+"Authorization header ausente", nil)
		utils.ErrorResponse(w, nil, http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		h.logger.Warn(ctx, ref+"Authorization header mal formatado", map[string]any{
			"header": authHeader,
		})
		utils.ErrorResponse(w, nil, http.StatusUnauthorized)
		return
	}

	tokenString := parts[1]

	err := h.service.Logout(ctx, tokenString)
	if err != nil {
		h.logger.Error(ctx, err, ref+"Erro ao invalidar token na blacklist", nil)
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+"Logout realizado com sucesso", nil)

	utils.ToJSON(w, http.StatusOK, utils.DefaultResponse{
		Status:  http.StatusOK,
		Message: "Logout realizado com sucesso",
	})
}
