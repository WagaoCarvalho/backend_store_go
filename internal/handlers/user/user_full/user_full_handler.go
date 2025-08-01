package handlers

import (
	"fmt"
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_full"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/users/user_full_services"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
)

type UserHandler struct {
	service services.UserFullService
	logger  *logger.LoggerAdapter
}

func NewUserFullHandler(service services.UserFullService, logger *logger.LoggerAdapter) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) CreateFull(w http.ResponseWriter, r *http.Request) {
	ref := "[UserHandler - CreateFull] "
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var requestData models.UserFull

	if err := utils.FromJson(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJsonError, map[string]any{
			"erro": err.Error(),
		})
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	createdUserFull, err := h.service.CreateFull(ctx, &requestData)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"email": requestData.User.Email,
		})
		utils.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":  createdUserFull.User.UID,
		"username": createdUserFull.User.Username,
		"email":    createdUserFull.User.Email,
	})

	utils.ToJson(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Usuário criado com sucesso",
		Data:    createdUserFull,
	})
}
