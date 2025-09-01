package handler

import (
	"fmt"
	"net/http"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user_full_services"
)

type UserHandler struct {
	service service.UserFullService
	logger  *logger.LogAdapter
}

func NewUserFullHandler(service service.UserFullService, logger *logger.LogAdapter) *UserHandler {
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

	if err := utils.FromJSON(r.Body, &requestData); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{
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

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Usuário criado com sucesso",
		Data:    createdUserFull,
	})
}
