package handler

import (
	"fmt"
	"net/http"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/full"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/full"
)

type UserHandler struct {
	service service.UserFull
	logger  *logger.LogAdapter
}

func NewUserFullHandler(service service.UserFull, logger *logger.LogAdapter) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) CreateFull(w http.ResponseWriter, r *http.Request) {
	const ref = "[UserHandler - CreateFull] "
	ctx := r.Context()

	if r.Method != http.MethodPost {
		h.logger.Warn(ctx, ref+logger.LogMethodNotAllowed, map[string]any{
			"method": r.Method,
		})
		utils.ErrorResponse(w, fmt.Errorf("método %s não permitido", r.Method), http.StatusMethodNotAllowed)
		return
	}

	h.logger.Info(ctx, ref+logger.LogCreateInit, nil)

	var requestDTO dto.UserFullDTO
	if err := utils.FromJSON(r.Body, &requestDTO); err != nil {
		h.logger.Warn(ctx, ref+logger.LogParseJSONError, map[string]any{"erro": err.Error()})
		utils.ErrorResponse(w, fmt.Errorf("erro ao decodificar JSON: %w", err), http.StatusBadRequest)
		return
	}

	if requestDTO.User == nil {
		h.logger.Warn(ctx, ref+logger.LogMissingBodyData, nil)
		utils.ErrorResponse(w, fmt.Errorf("dados do usuário são obrigatórios"), http.StatusBadRequest)
		return
	}

	// Converte DTO para model antes de passar para o service
	modelUserFull := dto.ToUserFullModel(requestDTO)

	createdUserFull, err := h.service.CreateFull(ctx, modelUserFull)
	if err != nil {
		h.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"email": modelUserFull.User.Email,
		})
		utils.ErrorResponse(w, fmt.Errorf("erro ao criar usuário: %w", err), http.StatusInternalServerError)
		return
	}

	// Converte model de volta para DTO para retornar
	createdDTO := dto.ToUserFullDTO(createdUserFull)

	h.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":  createdDTO.User.UID,
		"username": createdDTO.User.Username,
		"email":    createdDTO.User.Email,
	})

	utils.ToJSON(w, http.StatusCreated, utils.DefaultResponse{
		Status:  http.StatusCreated,
		Message: "Usuário criado com sucesso",
		Data:    createdDTO,
	})
}
